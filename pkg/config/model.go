package config

// Settings represents the configuration settings for the gh-dxp extension.
type Settings struct {
	Lint LintSettings `yaml:"lint"`
}

// LintSettings represents the configuration settings for the linting process.
type LintSettings struct {
	Linters []LinterSettings `yaml:"linters"`
	Exclude []string         `yaml:"exclude"`
}

// LinterSettings represents the configuration settings for a specific linter.
type LinterSettings struct {
	Name    string   `yaml:"name"`
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}
