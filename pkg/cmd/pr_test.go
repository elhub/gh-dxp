package cmd

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/spf13/cobra"
)

func TestAddPrOptionsToCreateOptions(t *testing.T) {
	prOptions := pr.Options{
		NoLint:        true,
		NoUnit:        true,
		CommitMessage: "test commit",
	}
	createOptions := &pr.CreateOptions{}

	addPrOptionsToCreateOptions(prOptions, createOptions)

	if createOptions.NoLint != prOptions.NoLint {
		t.Errorf("Expected NoLint to be %v, got %v", prOptions.NoLint, createOptions.NoLint)
	}
	if createOptions.NoUnit != prOptions.NoUnit {
		t.Errorf("Expected NoUnit to be %v, got %v", prOptions.NoUnit, createOptions.NoUnit)
	}
	if createOptions.CommitMessage != prOptions.CommitMessage {
		t.Errorf("Expected CommitMessage to be %v, got %v", prOptions.CommitMessage, createOptions.CommitMessage)
	}
}

func TestAddPrOptionsToUpdateOptions(t *testing.T) {
	prOptions := pr.Options{
		NoLint:        false,
		NoUnit:        true,
		CommitMessage: "update commit",
	}
	updateOptions := &pr.UpdateOptions{}

	addPrOptionsToUpdateOptions(prOptions, updateOptions)

	if updateOptions.NoLint != prOptions.NoLint {
		t.Errorf("Expected NoLint to be %v, got %v", prOptions.NoLint, updateOptions.NoLint)
	}
	if updateOptions.NoUnit != prOptions.NoUnit {
		t.Errorf("Expected NoUnit to be %v, got %v", prOptions.NoUnit, updateOptions.NoUnit)
	}
	if updateOptions.CommitMessage != prOptions.CommitMessage {
		t.Errorf("Expected CommitMessage to be %v, got %v", prOptions.CommitMessage, updateOptions.CommitMessage)
	}
}

func TestGetPrOptionsFromCmd(t *testing.T) {
	tests := []struct {
		name          string
		noLint        bool
		noUnit        bool
		commitMessage string
	}{
		{
			name:          "all flags set",
			noLint:        true,
			noUnit:        true,
			commitMessage: "test message",
		},
		{
			name:          "no flags set",
			noLint:        false,
			noUnit:        false,
			commitMessage: "",
		},
		{
			name:          "mixed flags",
			noLint:        true,
			noUnit:        false,
			commitMessage: "partial message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("nolint", tt.noLint, "")
			cmd.Flags().Bool("nounit", tt.noUnit, "")
			cmd.Flags().String("commitmessage", tt.commitMessage, "")

			result, err := getPrOptionsFromCmd(cmd)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result.NoLint != tt.noLint {
				t.Errorf("Expected NoLint to be %v, got %v", tt.noLint, result.NoLint)
			}
			if result.NoUnit != tt.noUnit {
				t.Errorf("Expected NoUnit to be %v, got %v", tt.noUnit, result.NoUnit)
			}
			if result.CommitMessage != tt.commitMessage {
				t.Errorf("Expected CommitMessage to be %v, got %v", tt.commitMessage, result.CommitMessage)
			}
		})
	}
}

func TestGetPrOptionsFromCmd_MissingFlags(t *testing.T) {
	cmd := &cobra.Command{}
	// Only add one flag to test error handling
	cmd.Flags().Bool("nolint", false, "")

	_, err := getPrOptionsFromCmd(cmd)

	if err == nil {
		t.Error("Expected error when flags are missing, got nil")
	}
}
