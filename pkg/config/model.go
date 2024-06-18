package config

// Settings represents the configuration settings for the gh-dxp extension.
type Settings struct {
	ProjectTemplateUri string `yaml:"projectTemplateUri"`
	ProjectType        string `yaml:"projectType"`
}
