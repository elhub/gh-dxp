// Package config defines configuration models for the gh-dxp extension.
package config

// Settings represents the configuration settings for the gh-dxp extension.
type Settings struct {
	ProjectTemplateURI     string `yaml:"projectTemplateUri"`
	ProjectType            string `yaml:"projectType"`
	JiraURL                string `yaml:"jiraUrl"`
	JiraEmail              string `yaml:"jiraEmail"`
	MegalinterImageVersion string `yaml:"megalinterImageVersion"`
}
