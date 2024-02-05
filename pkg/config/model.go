package config

type Settings struct {
	Lint LintSettings `yaml:"lint"`
}

type LintSettings struct {
	Linters []LinterSettings `yaml:"linters"`
	Exclude []string         `yaml:"exclude"`
}

type LinterSettings struct {
	Name    string   `yaml:"name"`
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}
