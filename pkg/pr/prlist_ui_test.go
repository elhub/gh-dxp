package pr

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialModel(t *testing.T) {
	testPRs := []PullRequestInfo{
		{
			Number: 123,
			Title:  "Test PR",
			Author: prAuthor{
				Login: "testuser",
				Name:  "Test User",
			},
			HeadRepository: prRepository{Name: "test-repo"},
			ReviewDecision: "APPROVED",
			Additions:      10,
			Deletions:      5,
		},
	}

	model := initialModel(testPRs)

	rows := model.table.Rows()
	if len(rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(rows))
	}

	if rows[0][0] != "test-repo" {
		t.Errorf("Expected repo 'test-repo', got '%s'", rows[0][0])
	}
	if rows[0][1] != "#123" {
		t.Errorf("Expected PR '#123', got '%s'", rows[0][1])
	}
	if rows[0][2] != "Test PR" {
		t.Errorf("Expected title 'Test PR', got '%s'", rows[0][2])
	}
	if rows[0][3] != "Test User" {
		t.Errorf("Expected author 'Test User', got '%s'", rows[0][3])
	}
	if rows[0][4] != "APPROVED" {
		t.Errorf("Expected status 'APPROVED', got '%s'", rows[0][4])
	}
	if rows[0][5] != "+10 -5" {
		t.Errorf("Expected changes '+10 -5', got '%s'", rows[0][5])
	}
}

func TestInitialModelEmpty(t *testing.T) {
	model := initialModel([]PullRequestInfo{})
	rows := model.table.Rows()

	if len(rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(rows))
	}
}

func TestInitialModelMultiple(t *testing.T) {
	testPRs := []PullRequestInfo{
		{
			Number:         1,
			Title:          "PR 1",
			Author:         prAuthor{Name: "User1"},
			HeadRepository: prRepository{Name: "repo1"},
			ReviewDecision: "APPROVED",
			Additions:      10,
			Deletions:      5,
		},
		{
			Number:         2,
			Title:          "PR 2",
			Author:         prAuthor{Name: "User2"},
			HeadRepository: prRepository{Name: "repo2"},
			ReviewDecision: "CHANGES_REQUESTED",
			Additions:      20,
			Deletions:      15,
		},
	}

	model := initialModel(testPRs)
	rows := model.table.Rows()

	if len(rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(rows))
	}

	if rows[1][0] != "repo2" {
		t.Errorf("Expected second repo 'repo2', got '%s'", rows[1][0])
	}
	if rows[1][5] != "+20 -15" {
		t.Errorf("Expected second PR changes '+20 -15', got '%s'", rows[1][5])
	}
}

func TestUIInit(t *testing.T) {
	model := initialModel([]PullRequestInfo{})
	cmd := model.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil")
	}
}

func TestUIUpdate(t *testing.T) {
	model := initialModel([]PullRequestInfo{})
	newModel, cmd := model.Update(tea.KeyMsg{})

	if cmd == nil {
		t.Error("Expected Update to return tea.Quit command")
	}

	_, ok := newModel.(PullRequestUI)
	if !ok {
		t.Error("Expected Update to return PullRequestUI model")
	}
}

func TestUIView(t *testing.T) {
	testPRs := []PullRequestInfo{
		{
			Number:         123,
			Title:          "Test",
			Author:         prAuthor{Name: "User"},
			HeadRepository: prRepository{Name: "repo"},
			ReviewDecision: "APPROVED",
			Additions:      10,
			Deletions:      5,
		},
	}

	model := initialModel(testPRs)
	view := model.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}
	if !strings.Contains(view, "Repository") {
		t.Error("Expected 'Repository' header in view")
	}
	if !strings.Contains(view, "Title") {
		t.Error("Expected 'Title' header in view")
	}
	if !strings.Contains(view, "repo") {
		t.Error("Expected repo name in view")
	}
	if !strings.Contains(view, "#123") {
		t.Error("Expected PR number in view")
	}
}

func TestChangesFormatting(t *testing.T) {
	tests := []struct {
		additions int
		deletions int
		expected  string
	}{
		{10, 5, "+10 -5"},
		{0, 0, "+0 -0"},
		{100, 200, "+100 -200"},
		{999, 1, "+999 -1"},
	}

	for _, tt := range tests {
		pr := PullRequestInfo{
			Number:         1,
			Title:          "Test",
			Author:         prAuthor{Name: "User"},
			HeadRepository: prRepository{Name: "repo"},
			ReviewDecision: "APPROVED",
			Additions:      tt.additions,
			Deletions:      tt.deletions,
		}

		model := initialModel([]PullRequestInfo{pr})
		rows := model.table.Rows()

		if rows[0][5] != tt.expected {
			t.Errorf("Expected changes '%s', got '%s'", tt.expected, rows[0][5])
		}
	}
}
