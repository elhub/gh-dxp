// Package projecttemplate provides utilities to set up new repositories using our project template.
package projecttemplate

// AskToConfirmFunc is a function type for confirming file deletion in tests.
type AskToConfirmFunc func(prompt string) (bool, error)

// Options represents the options for setting up a repository from the project template.
type Options struct {
	IsGradleProject bool
	TestRun         bool
	// CustomAskToConfirmFunc allows injecting a custom confirmation function for testing.
	// If nil, ghutil.AskToConfirm is used.
	CustomAskToConfirmFunc AskToConfirmFunc
}
