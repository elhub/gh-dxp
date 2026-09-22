// Package jira provides Jira issue search helpers used by gh-dxp commands.
package jira

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ErrJiraDisabled is returned when the user has opted out of Jira.
var ErrJiraDisabled = fmt.Errorf("jira disabled")

// ErrJiraNotConfigured is returned when ~/.jira_token does not exist and no env var is set.
var ErrJiraNotConfigured = fmt.Errorf("jira not configured")

// resolveCredentials reads Jira email and token with the following priority:
//  1. JIRA_USERNAME / JIRA_API_TOKEN env vars
//  2. ~/.jira_token file — expected format: "email:token" on a single line
//
// Returns ErrJiraDisabled if ~/.jira_token exists but is empty or contains "disabled" (user opted out).
// Returns ErrJiraNotConfigured if neither source has credentials.
func resolveCredentials(email string) (resolvedEmail, token string, err error) {
	resolvedEmail = email

	envToken := strings.TrimSpace(os.Getenv("JIRA_API_TOKEN"))
	if envEmail := strings.TrimSpace(os.Getenv("JIRA_USERNAME")); envEmail != "" {
		resolvedEmail = envEmail
	}
	if envToken != "" {
		return resolvedEmail, envToken, nil
	}

	tokenFile := filepath.Join(os.Getenv("HOME"), ".jira_token")
	data, fileErr := os.ReadFile(tokenFile)
	if os.IsNotExist(fileErr) {
		return "", "", ErrJiraNotConfigured
	}
	if fileErr != nil {
		return "", "", fileErr
	}

	content := strings.TrimSpace(string(data))
	if content == "" || content == "disabled" {
		return "", "", ErrJiraDisabled
	}

	// Format: email:token
	parts := strings.SplitN(content, ":", 2)
	if len(parts) == 2 {
		resolvedEmail = strings.TrimSpace(parts[0])
		token = strings.TrimSpace(parts[1])
	} else {
		token = content
	}
	return resolvedEmail, token, nil
}

// SearchText contains text used to rank Jira issue search results.
type SearchText struct {
	CommitMessage string
	Title         string
	Description   string
}

// SearchIssue represents a Jira issue returned by the search API.
type SearchIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string          `json:"summary"`
		Description json.RawMessage `json:"description"`
	} `json:"fields"`
}

// SearchIssues searches Jira issues and returns results ranked by relevance.
func SearchIssues(ctx context.Context, baseURL, email string, text SearchText) ([]SearchIssue, error) {
	baseURL = strings.TrimRight(strings.TrimSuffix(baseURL, "/browse"), "/")

	resolvedEmail, token, err := resolveCredentials(email)
	if err != nil {
		return nil, err
	}

	if baseURL == "" || resolvedEmail == "" {
		return nil, fmt.Errorf("Jira URL and email are required")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	query := url.Values{
		"jql":        {"assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC"},
		"fields":     {"summary,description"},
		"maxResults": {"50"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/rest/api/3/search/jql?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(resolvedEmail+":"+token)))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Jira request failed with status %s", resp.Status)
	}

	var result struct {
		Issues []SearchIssue `json:"issues"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return rankIssues(result.Issues, text), nil
}

func rankIssues(issues []SearchIssue, text SearchText) []SearchIssue {
	type scored struct {
		issue SearchIssue
		score int
	}
	var matches []scored
	for _, issue := range issues {
		summary := strings.ToLower(issue.Fields.Summary)
		description := strings.ToLower(string(issue.Fields.Description))
		score := scoreText(summary, text.CommitMessage, 6) + scoreText(description, text.CommitMessage, 3)
		score += scoreText(summary, text.Title, 3) + scoreText(description, text.Title, 2)
		score += scoreText(summary, text.Description, 1) + scoreText(description, text.Description, 1)
		if score > 0 {
			matches = append(matches, scored{issue, score})
		}
	}
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].score > matches[j].score })
	result := make([]SearchIssue, len(matches))
	for i := range matches {
		result[i] = matches[i].issue
	}
	return result
}

func scoreText(content, text string, weight int) int {
	score := 0
	for _, term := range strings.Fields(strings.ToLower(text)) {
		if strings.Contains(content, term) {
			score += weight
		}
	}
	return score
}
