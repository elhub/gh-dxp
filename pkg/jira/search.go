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
	"unicode"
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
		if resolvedEmail == "" {
			return "", "", ErrJiraNotConfigured
		}
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

// isTrustedJiraURL allows only the Elhub Jira endpoint. It is a variable so tests can use a local server.
var isTrustedJiraURL = func(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "https" {
		return false
	}
	return u.Host == "elhub.atlassian.net" && u.User == nil &&
		u.Path == "" && u.RawQuery == "" && u.Fragment == ""
}

// SearchIssues searches Jira issues and returns results ranked by relevance.
func SearchIssues(ctx context.Context, baseURL, email string, text SearchText) ([]SearchIssue, error) {
	baseURL = strings.TrimSuffix(strings.TrimRight(baseURL, "/"), "/browse")

	resolvedEmail, token, err := resolveCredentials(email)
	if err != nil {
		return nil, err
	}

	if baseURL == "" || resolvedEmail == "" {
		return nil, fmt.Errorf("Jira URL and email are required")
	}
	// Only send credentials to the Elhub Jira endpoint.
	if !isTrustedJiraURL(baseURL) {
		return nil, fmt.Errorf("refusing to send Jira credentials to untrusted URL %q", baseURL)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	query := url.Values{
		"jql":        {"assignee = currentUser() AND statusCategory != Done AND project != ET ORDER BY updated DESC"},
		"fields":     {"summary,description"},
		"maxResults": {"50"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/rest/api/3/search/jql?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(resolvedEmail+":"+token)))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
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
		summary := tokenize(issue.Fields.Summary)
		description := tokenize(string(issue.Fields.Description))
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

// stopWords are ignored when matching. It includes the structural keys of Jira's JSON description format.
var stopWords = map[string]bool{
	"the": true, "and": true, "for": true, "with": true, "from": true, "this": true, "that": true,
	"add": true, "fix": true, "feat": true, "chore": true, "docs": true, "test": true, "refactor": true,
	"type": true, "text": true, "content": true, "doc": true, "paragraph": true, "version": true,
}

// tokenize lowercases s and splits it into words, dropping short words and stop words.
func tokenize(s string) map[string]bool {
	words := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	tokens := make(map[string]bool, len(words))
	for _, w := range words {
		if len(w) >= 3 && !stopWords[w] {
			tokens[w] = true
		}
	}
	return tokens
}

// scoreText adds weight for every distinct term in text that is a whole word of content.
func scoreText(content map[string]bool, text string, weight int) int {
	score := 0
	for term := range tokenize(text) {
		if content[term] {
			score += weight
		}
	}
	return score
}
