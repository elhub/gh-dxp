package opslevel

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestExecute(t *testing.T) {
	dir := t.TempDir()
	exe := new(testutils.MockExecutor)
	exe.On("GH", []string{"repo", "view", "--json", "name,nameWithOwner,description,url,primaryLanguage"}).Return(
		`{"name":"my-tool","nameWithOwner":"elhub/my-tool","description":"A tool: for developers.",`+
			`"url":"https://github.com/elhub/my-tool","primaryLanguage":{"name":"Go"}}`, nil)
	options := Options{
		Type: "tool", Lifecycle: "production", Tier: "tier_4", Owner: "team-devxp",
		PackageName: "no.elhub.my-tool", NoInput: true,
	}
	require.NoError(t, Execute(dir, exe, options))
	content, err := os.ReadFile(filepath.Join(dir, "opslevel.yml"))
	require.NoError(t, err)
	var document struct {
		Version   int       `yaml:"version"`
		Component component `yaml:"component"`
	}
	require.NoError(t, yaml.Unmarshal(content, &document))
	require.Equal(t, 1, document.Version)
	require.Equal(t, "my-tool", document.Component.Name)
	require.Equal(t, "A tool: for developers.", document.Component.Description)
	require.Equal(t, "Go", document.Component.Language)
	require.Equal(t, []string{"my_tool", "no.elhub.my-tool"}, document.Component.Aliases)
	require.Equal(t, "elhub/my-tool", document.Component.Repositories[0].Name)
	require.Equal(t, "https://github.com/elhub/my-tool", document.Component.Tools[0].URL)
	require.ErrorContains(t, Execute(dir, exe, options), "refusing to overwrite")
	after, err := os.ReadFile(filepath.Join(dir, "opslevel.yml"))
	require.NoError(t, err)
	require.Equal(t, content, after)
	exe.AssertNumberOfCalls(t, "GH", 1)
	exe.AssertExpectations(t)
}

func TestCompleteOptions(t *testing.T) {
	options := Options{}
	choiceAnswers := []string{"tool", "development", "tier_3", "None", "None", "None"}
	stringAnswers := []string{"team-devxp", "Development tool", "no.elhub.my-tool"}
	err := completeOptions(&options, func(_ string, _ string) (string, error) {
		answer := stringAnswers[0]
		stringAnswers = stringAnswers[1:]
		return answer, nil
	}, func(_ string, choices []string) (string, error) {
		answer := choiceAnswers[0]
		choiceAnswers = choiceAnswers[1:]
		require.Contains(t, choices, answer)
		return answer, nil
	})
	require.NoError(t, err)
	require.Equal(t, "team-devxp", options.Owner)
	require.Empty(t, choiceAnswers)
	require.Empty(t, stringAnswers)
}

func TestCompleteOptionsRejectsInvalidValues(t *testing.T) {
	for _, field := range []string{"type", "lifecycle", "tier", "owner", "description", "package-name"} {
		t.Run(field, func(t *testing.T) {
			options := Options{Type: "tool", Lifecycle: "production", Tier: "tier_4", Owner: "team-devxp", Description: "A tool", NoInput: true}
			switch field {
			case "type":
				options.Type = "unknown"
			case "lifecycle":
				options.Lifecycle = "active"
			case "tier":
				options.Tier = "4"
			case "owner":
				options.Owner = " "
			case "description":
				options.Description = ""
			case "package-name":
				options.PackageName = " "
			}
			require.ErrorContains(t, completeOptions(&options, nil, nil), "--"+field)
		})
	}
}

func TestExecuteFailuresDoNotCreateFile(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		metadata string
		err      error
		want     string
	}{
		{"github error", "", errors.New("not authenticated"), "could not read repository metadata"},
		{"invalid json", "not json", nil, "could not parse repository metadata"},
		{"missing repository", `{}`, nil, "repository metadata is missing"},
		{
			"missing owner",
			`{"name":"example","nameWithOwner":"elhub/example","url":"https://github.com/elhub/example"}`,
			nil, "--owner is required",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			dir := t.TempDir()
			exe := new(testutils.MockExecutor)
			exe.On("GH", []string{"repo", "view", "--json", "name,nameWithOwner,description,url,primaryLanguage"}).Return(
				testCase.metadata, testCase.err)
			options := Options{Type: "tool", Lifecycle: "development", Tier: "tier_4", NoInput: true}
			require.ErrorContains(t, Execute(dir, exe, options), testCase.want)
			require.NoFileExists(t, filepath.Join(dir, "opslevel.yml"))
			exe.AssertExpectations(t)
		})
	}
}

func TestCompleteOptionsCancellation(t *testing.T) {
	cancelled := errors.New("prompt cancelled")
	options := Options{}
	err := completeOptions(&options, nil, func(_ string, _ []string) (string, error) {
		return "", cancelled
	})
	require.ErrorIs(t, err, cancelled)
	options = Options{Type: "tool", Lifecycle: "development", Tier: "tier_4"}
	err = completeOptions(&options, func(_ string, _ string) (string, error) {
		return "", cancelled
	}, func(_ string, choices []string) (string, error) {
		return choices[0], nil
	})
	require.ErrorIs(t, err, cancelled)
}

func TestReviewBeforeWriting(t *testing.T) {
	for _, updating := range []bool{false, true} {
		for _, decision := range []string{"accept", "decline", "cancel"} {
			name := "create/" + decision
			if updating {
				name = "update/" + decision
			}
			t.Run(name, func(t *testing.T) {
				dir := t.TempDir()
				path := filepath.Join(dir, "opslevel.yml")
				original := "version: 1\ncomponent:\n  name: my-tool\n  type: tool\n" +
					"  lifecycle: development\n  tier: tier_3\n  owner: team-existing\n" +
					"  description: Existing description\n  tags: [{key: datastore, value: postgresql}]\n"
				options := Options{Type: "tool", Lifecycle: "production", Tier: "tier_4", Owner: "team-devxp", PackageName: "no.elhub.my-tool"}
				if updating {
					require.NoError(t, os.WriteFile(path, []byte(original), 0o600))
					options = Options{Lifecycle: "production", PackageName: "no.elhub.my-tool"}
				}
				exe := new(testutils.MockExecutor)
				exe.On("GH", []string{"repo", "view", "--json", "name,nameWithOwner,description,url,primaryLanguage"}).Return(
					`{"name":"my-tool","nameWithOwner":"elhub/my-tool","description":"GitHub description","url":"https://github.com/elhub/my-tool"}`, nil)
				reviews := 0
				var preview string
				cancelled := errors.New("cancelled")
				err := execute(dir, exe, options, func(prompt string) (bool, error) {
					if strings.Contains(prompt, "already exists") {
						return true, nil
					}
					reviews++
					require.Contains(t, prompt, "Write these values?")
					preview = strings.TrimSuffix(strings.TrimPrefix(prompt, "Review opslevel.yml:\n\n"), "\nWrite these values?")
					require.Contains(t, preview, "lifecycle: production")
					if updating {
						require.Contains(t, preview, "owner: team-existing")
						require.Contains(t, preview, "description: |-\n    Existing description")
						require.Contains(t, preview, "value: postgresql")
					}
					if decision == "cancel" {
						return false, cancelled
					}
					return decision == "accept", nil
				}, func(_ string, defaultValue string) (string, error) {
					return defaultValue, nil
				}, func(_ string, choices []string) (string, error) {
					return choices[0], nil
				})
				require.Equal(t, 1, reviews)
				if decision == "accept" {
					require.NoError(t, err)
					content, err := os.ReadFile(path)
					require.NoError(t, err)
					require.Equal(t, preview, string(content))
				} else {
					if decision == "cancel" {
						require.ErrorIs(t, err, cancelled)
					} else {
						require.ErrorIs(t, err, ErrWriteDeclined)
					}
					if updating {
						content, err := os.ReadFile(path)
						require.NoError(t, err)
						require.Equal(t, original, string(content))
					} else {
						require.NoFileExists(t, path)
					}
				}
				exe.AssertExpectations(t)
			})
		}
	}
}

func TestCompleteOptionsReviewsDefaults(t *testing.T) {
	options := Options{
		Type: "tool", Lifecycle: "production", Tier: "tier_4", Owner: "team-devxp",
		System: "devxp", Language: "Go", Framework: "Cobra",
		Description: "Original description", PackageName: "no.elhub.my-tool",
	}
	choiceDefaults := []string{"tool", "production", "tier_4", "devxp", "Go", "Cobra"}
	stringDefaults := []string{"team-devxp", "Original description", "no.elhub.my-tool"}
	err := completeOptions(&options, func(_ string, defaultValue string) (string, error) {
		require.NotEmpty(t, stringDefaults)
		require.Equal(t, stringDefaults[0], defaultValue)
		stringDefaults = stringDefaults[1:]
		return defaultValue + " updated", nil
	}, func(prompt string, choices []string) (string, error) {
		require.NotEmpty(t, choiceDefaults)
		require.Equal(t, choiceDefaults[0], choices[0])
		require.Contains(t, prompt, "default: "+choiceDefaults[0])
		choiceDefaults = choiceDefaults[1:]
		if strings.HasPrefix(prompt, "System") {
			return "developer_experience", nil
		}
		if strings.HasPrefix(prompt, "Language") {
			return "Kotlin", nil
		}
		if strings.HasPrefix(prompt, "Framework") {
			return "Ktor", nil
		}
		return choices[1], nil
	})
	require.NoError(t, err)
	require.Empty(t, choiceDefaults)
	require.Empty(t, stringDefaults)
	require.NotEqual(t, "tool", options.Type)
	require.NotEqual(t, "production", options.Lifecycle)
	require.NotEqual(t, "tier_4", options.Tier)
	require.Equal(t, "team-devxp updated", options.Owner)
	require.Equal(t, "developer_experience", options.System)
	require.Equal(t, "Kotlin", options.Language)
	require.Equal(t, "Ktor", options.Framework)
	require.Equal(t, "Original description updated", options.Description)
}

func TestPackageNameRequiredBeforeWriting(t *testing.T) {
	dir := t.TempDir()
	exe := new(testutils.MockExecutor)
	exe.On("GH", []string{"repo", "view", "--json", "name,nameWithOwner,description,url,primaryLanguage"}).Return(
		`{"name":"my-tool","nameWithOwner":"elhub/my-tool","description":"A tool","url":"https://github.com/elhub/my-tool"}`, nil)
	options := Options{Type: "tool", Lifecycle: "production", Tier: "tier_4", Owner: "team-devxp", NoInput: true}
	require.ErrorContains(t, Execute(dir, exe, options), "--package-name is required")
	require.NoFileExists(t, filepath.Join(dir, "opslevel.yml"))
}

func TestRequiredAliasesMergedOnUpdate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opslevel.yml")
	require.NoError(t, os.WriteFile(path, []byte("version: 1\ncomponent:\n  aliases: [custom, my_tool, no.elhub.my-tool, custom]\n"), 0o600))
	exe := new(testutils.MockExecutor)
	exe.On("GH", []string{"repo", "view", "--json", "name,nameWithOwner,description,url,primaryLanguage"}).Return(
		`{"name":"my-tool","nameWithOwner":"elhub/my-tool","description":"A tool","url":"https://github.com/elhub/my-tool"}`, nil)
	options := Options{Type: "tool", Lifecycle: "production", Tier: "tier_4", Owner: "team-devxp", PackageName: "no.elhub.my-tool"}
	err := execute(dir, exe, options, func(_ string) (bool, error) {
		return true, nil
	}, func(_ string, defaultValue string) (string, error) {
		return defaultValue, nil
	}, func(_ string, choices []string) (string, error) {
		return choices[0], nil
	})
	require.NoError(t, err)
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	var document struct {
		Component component `yaml:"component"`
	}
	require.NoError(t, yaml.Unmarshal(content, &document))
	require.Equal(t, []string{"my_tool", "no.elhub.my-tool", "custom"}, document.Component.Aliases)
}

func TestSystemChoices(t *testing.T) {
	for _, testCase := range []struct {
		name, current, selection, want string
	}{
		{"known", "", "settlement", "settlement"},
		{"existing known", "market_data", "market_data", "market_data"},
		{"existing custom", "my-system", "my-system", "my-system"},
		{"custom input", "market_data", "Custom...", "custom-system"},
		{"clear", "market_data", "None", ""},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			textPrompted := false
			value, err := askSystem(testCase.current, func(_ string, defaultValue string) (string, error) {
				textPrompted = true
				require.Equal(t, testCase.current, defaultValue)
				return "custom-system", nil
			}, func(prompt string, choices []string) (string, error) {
				defaultValue := testCase.current
				if defaultValue == "" {
					defaultValue = "None"
				}
				require.Equal(t, defaultValue, choices[0])
				require.Contains(t, prompt, "default: "+defaultValue)
				for _, system := range []string{
					"core_platform", "developer_experience", "documentation_platform", "end_user_portal",
					"identity_and_access_management", "market_agent", "market_data", "market_portal",
					"metering_values", "service_platform", "settlement", "None", "Custom...",
				} {
					require.Contains(t, choices, system)
				}
				return testCase.selection, nil
			})
			require.NoError(t, err)
			require.Equal(t, testCase.want, value)
			require.Equal(t, testCase.selection == "Custom...", textPrompted)
		})
	}
}

func TestLanguageChoices(t *testing.T) {
	for _, selection := range []string{"Kotlin", "Java", "Python", "Go", "Typescript", "Rust", "None", "Custom..."} {
		t.Run(selection, func(t *testing.T) {
			textPrompted := false
			value, err := askLanguage("Rust", func(prompt, defaultValue string) (string, error) {
				textPrompted = true
				require.Equal(t, "Custom language:", prompt)
				require.Equal(t, "Rust", defaultValue)
				return "Scala", nil
			}, func(prompt string, choices []string) (string, error) {
				require.Contains(t, prompt, "default: Rust")
				require.Equal(t, []string{"Rust", "None", "Kotlin", "Java", "Python", "Go", "Typescript", "Custom..."}, choices)
				return selection, nil
			})
			require.NoError(t, err)
			want := selection
			if selection == "None" {
				want = ""
			} else if selection == "Custom..." {
				want = "Scala"
			}
			require.Equal(t, want, value)
			require.Equal(t, selection == "Custom...", textPrompted)
		})
	}
}

func TestFrameworkChoices(t *testing.T) {
	for _, selection := range []string{"Ktor", "React", "Ansible", "Cobra", "None", "Custom..."} {
		t.Run(selection, func(t *testing.T) {
			textPrompted := false
			value, err := askFramework("Cobra", func(prompt, defaultValue string) (string, error) {
				textPrompted = true
				require.Equal(t, "Custom framework:", prompt)
				require.Equal(t, "Cobra", defaultValue)
				return "Spring", nil
			}, func(prompt string, choices []string) (string, error) {
				require.Contains(t, prompt, "default: Cobra")
				require.Equal(t, []string{"Cobra", "None", "Ktor", "React", "Ansible", "Custom..."}, choices)
				return selection, nil
			})
			require.NoError(t, err)
			want := selection
			if selection == "None" {
				want = ""
			} else if selection == "Custom..." {
				want = "Spring"
			}
			require.Equal(t, want, value)
			require.Equal(t, selection == "Custom...", textPrompted)
		})
	}
}

func TestDescriptionBlockRoundTrip(t *testing.T) {
	for _, description := range []string{"Single paragraph.", "First line.\nSecond line.\n\nAnother paragraph."} {
		content, err := orderedYAML(map[string]any{"version": 1, "component": map[string]any{"description": description}})
		require.NoError(t, err)
		require.Contains(t, string(content), "description: |-\n")
		var document struct {
			Component component `yaml:"component"`
		}
		require.NoError(t, yaml.Unmarshal(content, &document))
		require.Equal(t, description, document.Component.Description)
	}
}

func TestReferenceFieldOrder(t *testing.T) {
	componentOrder := []string{
		"name", "type", "lifecycle", "tier", "owner", "system", "language", "framework", "description",
		"aliases", "tags", "repositories", "tools", "dependencies", "alert_sources", "properties", "custom",
	}
	fields := map[string]any{}
	for _, key := range componentOrder {
		fields[key] = "value"
	}
	entryOrders := map[string][]string{
		"tags": {"key", "value"}, "repositories": {"name", "path", "provider", "display_name"},
		"tools": {"name", "category", "url", "environment"}, "dependencies": {"alias", "notes"},
	}
	for key, order := range entryOrders {
		entry := map[string]any{}
		for _, field := range order {
			entry[field] = "value"
		}
		fields[key] = []any{entry}
	}
	content, err := orderedYAML(map[string]any{"component": fields, "version": 1})
	require.NoError(t, err)
	var document yaml.Node
	require.NoError(t, yaml.Unmarshal(content, &document))
	root := document.Content[0]
	require.Equal(t, "version", root.Content[0].Value)
	require.Equal(t, "component", root.Content[2].Value)
	componentNode := root.Content[3]
	var keys []string
	for index := 0; index < len(componentNode.Content); index += 2 {
		key := componentNode.Content[index].Value
		keys = append(keys, key)
		if order, ok := entryOrders[key]; ok {
			entry := componentNode.Content[index+1].Content[0]
			var entryKeys []string
			for entryIndex := 0; entryIndex < len(entry.Content); entryIndex += 2 {
				entryKeys = append(entryKeys, entry.Content[entryIndex].Value)
			}
			require.Equal(t, order, entryKeys)
		}
	}
	require.Equal(t, componentOrder, keys)
}
