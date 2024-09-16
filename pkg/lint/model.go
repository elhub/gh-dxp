package lint

// Options represents the options for the lint command.
type Options struct {
	Fix         bool
	LintAll     bool
	Directory   string
	LinterImage string
}
