package status //nolint:cyclop

// Options represents the options available for the status command.
type Options struct {
	All      bool
	Repo     bool
	Pr       bool
	Branches bool
	Issue    bool
}
