// pkg/status/status.go

package status

import (
    "fmt"
    "strings"
    "github.com/elhub/gh-dxp/pkg/utils"
)

// Status encapsulates methods related to getting the status of the repository.
type Status struct {
    Executor utils.Executor // Use utils.Executor from the utils package
}

// NewStatus creates a new Status instance with the provided executor.
func NewStatus(executor utils.Executor) *Status {
    return &Status{Executor: executor}
}

// GetStatus retrieves the status of the current repository.
func (s *Status) GetStatus(statusType string) (string, error) {
    var statusReport strings.Builder

    if statusType == "All" || statusType == "Current repository" {
        repo, err := s.Executor.Command("git", "remote", "get-url", "origin")
        if err != nil {
            return "", err
        }
        statusReport.WriteString(fmt.Sprintf("Repository: %s\n", strings.TrimSpace(repo)))
    }

    if statusType == "All" || statusType == "PR status" {
        prStatus, err := s.Executor.Command("gh", "pr", "status")
        if err != nil {
            return "", err
        }
        statusReport.WriteString(fmt.Sprintf("PR Status:\n%s\n", prStatus))
    }

    if statusType == "All" || statusType == "List branches" {
        branches, err := s.Executor.Command("git", "branch", "-a")
        if err != nil {
            return "", err
        }
        statusReport.WriteString(fmt.Sprintf("Branches:\n%s\n", branches))
    }

    if statusType == "All" || statusType == "Assigned PRs/Review Requests" {
        assignedPRs, err := s.Executor.Command("gh", "issue", "status")
        if err != nil {
            return "", err
        }
        statusReport.WriteString(fmt.Sprintf("Assigned PRs/Review Requests:\n%s\n", assignedPRs))
    }

    return statusReport.String(), nil
}
