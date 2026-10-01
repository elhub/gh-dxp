package pr

import (
	"encoding/json"
	"testing"

	"github.com/elhub/gh-dxp/pkg/jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatJiraSuggestions(t *testing.T) {
	var issues []jira.SearchIssue
	require.NoError(t, json.Unmarshal([]byte(`[
		{"key":"TDX-1","fields":{"summary":"Best match"}},
		{"key":"TDX-2","fields":{"summary":"Second match"}},
		{"key":"TDX-3","fields":{"summary":"Third match"}},
		{"key":"TDX-4","fields":{"summary":"Fourth match"}},
		{"key":"TDX-5","fields":{"summary":"Fifth match"}},
		{"key":"TDX-6","fields":{"summary":"Not shown"}}
	]`), &issues))

	expected := "Suggested issues:\n" +
		"- TDX-1 - Best match\n" +
		"- TDX-2 - Second match\n" +
		"- TDX-3 - Third match\n" +
		"- TDX-4 - Fourth match\n" +
		"- TDX-5 - Fifth match"
	assert.Equal(t, expected, formatJiraSuggestions(issues))
	assert.Empty(t, formatJiraSuggestions(nil))
}

func TestJiraIssuePromptDefault(t *testing.T) {
	detectedIDs := []string{"ABC-10"}
	suggestions := []jira.SearchIssue{{Key: "ABC-20"}}

	assert.Equal(t, "ABC-20", jiraIssuePromptDefault(detectedIDs, suggestions))
	assert.Equal(t, "ABC-10", jiraIssuePromptDefault(detectedIDs, nil))
	assert.Equal(t, "ABC-10", jiraIssuePromptDefault(detectedIDs, []jira.SearchIssue{{Key: "ABC-20"}, {Key: "ABC-30"}}))
}
