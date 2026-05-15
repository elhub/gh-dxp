// Package projecttemplate provides utilities to set up new repositories using our project template.
package projecttemplate

// Options represents the options for setting up a repository from the project template.
type Options struct {
	IsGradleProject bool
	TestRun         bool
}
