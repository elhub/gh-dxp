package utils

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// ConvertTerminalOutputIntoList converts terminal output on multiple lines to a list of strings
func ConvertTerminalOutputIntoList(changedFilesString string) []string {
	if len(changedFilesString) == 0 {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(changedFilesString), "\n")
}

// AskToConfirm prompts the user with a yes/no question and returns their response.
func AskToConfirm(question string) (bool, error) {
	confirm := false
	err := survey.AskOne(&survey.Confirm{
		Message: question,
	}, &confirm, survey.WithValidator(survey.Required))
	if err != nil {
		return false, err
	}
	return confirm, nil
}

// AskForString prompts the user for a string input and returns the response.
func AskForString(question string, defaultAnswer string) (string, error) {
	var title string
	prompt := &survey.Input{
		Message: question,
		Default: defaultAnswer,
	}
	err := survey.AskOne(prompt, &title)
	if err != nil {
		return "", err
	}
	return title, nil
}

// AskForMultiline prompts the user for a multiline input and returns the response.
func AskForMultiline(question string) (string, error) {
	lines := ""
	err := survey.AskOne(&survey.Multiline{
		Message: question,
	}, &lines, survey.WithValidator(survey.Required))
	if err != nil {
		return lines, err
	}
	return lines, nil
}
