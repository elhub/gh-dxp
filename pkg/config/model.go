package config

// Settings represents the configuration settings for the gh-dxp extension.
type Settings struct {
	ProjectTemplateURI string `yaml:"projectTemplateUri"`
	ProjectType        string `yaml:"projectType"`
	TicketingBaseURL   string `yaml:"ticketingBaseUrl"`
}
