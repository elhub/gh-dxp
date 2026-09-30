// Package lint provides utilities for linting source code in gh-dxp.
package lint

// Options represents the options for the lint command.
type Options struct {
	Fix         bool
	LintAll     bool
	Directory   string
	LinterImage string
	Proxy       string
	// BaseDir controls where lint config files are discovered (for example
	// .mega-linter.yml). When set, Run passes MEGALINTER_CONFIG from this path
	// explicitly. It does not change the process working directory used to run
	// mega-linter-runner.
	BaseDir     string
}
