// Package opslevel generates and updates repository-based OpsLevel catalog configuration.
package opslevel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/elhub/gh-dxp/pkg/ghutil"
	"gopkg.in/yaml.v3"
)

// Options contains component metadata and controls interactive configuration prompts.
type Options struct {
	Type        string
	Lifecycle   string
	Tier        string
	Owner       string
	System      string
	Language    string
	Framework   string
	Description string
	PackageName string
	NoInput     bool
}

type component struct {
	Name         string       `yaml:"name"`
	Type         string       `yaml:"type"`
	Lifecycle    string       `yaml:"lifecycle"`
	Tier         string       `yaml:"tier"`
	Owner        string       `yaml:"owner"`
	System       string       `yaml:"system,omitempty"`
	Language     string       `yaml:"language,omitempty"`
	Framework    string       `yaml:"framework,omitempty"`
	Description  string       `yaml:"description"`
	Aliases      []string     `yaml:"aliases,omitempty"`
	Repositories []repository `yaml:"repositories"`
	Tools        []tool       `yaml:"tools"`
}

type repository struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Provider    string `yaml:"provider"`
	DisplayName string `yaml:"display_name"`
}

type tool struct {
	Name        string `yaml:"name"`
	Category    string `yaml:"category"`
	URL         string `yaml:"url"`
	Environment string `yaml:"environment"`
}

// ErrOverwriteDeclined indicates that the user declined to overwrite an existing configuration.
var ErrOverwriteDeclined = errors.New("opslevel.yml left unchanged")

// ErrWriteDeclined indicates that the user declined the final configuration review.
var ErrWriteDeclined = errors.New("no opslevel.yml changes written")

// Execute creates or updates opslevel.yml using repository metadata and user-provided values.
func Execute(workingDir string, exe ghutil.Executor, options Options) error {
	return execute(workingDir, exe, options, ghutil.AskToConfirm, ghutil.AskForString, ghutil.AskForMultipleChoice)
}

func execute(
	workingDir string,
	exe ghutil.Executor,
	options Options,
	confirm func(string) (bool, error),
	askString func(string, string) (string, error),
	askChoice func(string, []string) (string, error),
) error {
	path := filepath.Join(workingDir, "opslevel.yml")
	existing, mode, err := readExisting(path, options.NoInput, confirm)
	if err != nil {
		return err
	}
	if existing != nil {
		if err := applyExistingDefaults(&options, existing["component"]); err != nil {
			return err
		}
	}

	output, err := exe.GH("repo", "view", "--json", "name,nameWithOwner,description,url,primaryLanguage")
	if err != nil {
		return fmt.Errorf("could not read repository metadata: %w", err)
	}
	var metadata struct {
		Name            string
		NameWithOwner   string
		Description     string
		URL             string
		PrimaryLanguage *struct{ Name string }
	}
	if err := json.Unmarshal([]byte(output), &metadata); err != nil {
		return fmt.Errorf("could not parse repository metadata: %w", err)
	}
	if metadata.Name == "" || metadata.NameWithOwner == "" || metadata.URL == "" {
		return fmt.Errorf("repository metadata is missing name, nameWithOwner, or url")
	}
	if options.Description == "" {
		options.Description = metadata.Description
	}
	if options.Language == "" && metadata.PrimaryLanguage != nil {
		options.Language = metadata.PrimaryLanguage.Name
	}
	if err := completeOptions(&options, askString, askChoice); err != nil {
		return err
	}

	catalogComponent := component{
		Name: metadata.Name, Type: options.Type, Lifecycle: options.Lifecycle,
		Tier: options.Tier, Owner: options.Owner, System: options.System,
		Language: options.Language, Framework: options.Framework, Description: options.Description,
		Repositories: []repository{{Name: metadata.NameWithOwner, Path: "/", Provider: "github", DisplayName: metadata.Name}},
		Tools:        []tool{{Name: "Source Code", Category: "code", URL: metadata.URL, Environment: "Production"}},
	}
	catalogComponent.Aliases = []string{strings.ReplaceAll(metadata.Name, "-", "_")}
	if !slices.Contains(catalogComponent.Aliases, options.PackageName) {
		catalogComponent.Aliases = append(catalogComponent.Aliases, options.PackageName)
	}
	document := struct {
		Version   int       `yaml:"version"`
		Component component `yaml:"component"`
	}{Version: 1, Component: catalogComponent}
	content, err := orderedYAML(document)
	if err != nil {
		return err
	}
	if existing != nil {
		var generated map[string]any
		if err := yaml.Unmarshal(content, &generated); err != nil {
			return err
		}
		fields := existing["component"].(map[string]any)
		for _, key := range []string{"system", "language", "framework"} {
			delete(fields, key)
		}
		for key, value := range generated["component"].(map[string]any) {
			if key == "aliases" {
				aliases := value.([]any)
				if previous, ok := fields[key].([]any); ok {
					for _, alias := range previous {
						if !slices.Contains(aliases, alias) {
							aliases = append(aliases, alias)
						}
					}
				}
				value = aliases
			}
			if _, present := fields[key]; present && slices.Contains([]string{"repositories", "tools"}, key) {
				continue
			}
			fields[key] = value
		}
		content, err = orderedYAML(existing)
		if err != nil {
			return err
		}
	}
	if !options.NoInput {
		approved, err := confirm("Review opslevel.yml:\n\n---\n" + string(content) + "\nWrite these values?")
		if err != nil {
			return err
		}
		if !approved {
			return ErrWriteDeclined
		}
	}
	if existing != nil {
		return replaceExisting(path, append([]byte("---\n"), content...), mode)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("could not create opslevel.yml: %w", err)
	}
	_, writeErr := file.Write(append([]byte("---\n"), content...))
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

func orderedYAML(document any) ([]byte, error) {
	var root yaml.Node
	if err := root.Encode(document); err != nil {
		return nil, err
	}
	orderMapping(&root, []string{"version", "component"})
	for index := 0; index < len(root.Content); index += 2 {
		if root.Content[index].Value != "component" {
			continue
		}
		fields := root.Content[index+1]
		orderMapping(fields, []string{
			"name", "type", "lifecycle", "tier", "owner", "system", "language", "framework",
			"description", "aliases", "tags", "repositories", "tools", "dependencies", "alert_sources", "properties",
		})
		entryOrders := map[string][]string{
			"tags":         {"key", "value"},
			"repositories": {"name", "path", "provider", "display_name"},
			"tools":        {"name", "category", "url", "environment"},
			"dependencies": {"alias", "notes"},
		}
		for fieldIndex := 0; fieldIndex < len(fields.Content); fieldIndex += 2 {
			order, found := entryOrders[fields.Content[fieldIndex].Value]
			entries := fields.Content[fieldIndex+1]
			if fields.Content[fieldIndex].Value == "description" && entries.Kind == yaml.ScalarNode {
				entries.Style = yaml.LiteralStyle
			}
			if found && entries.Kind == yaml.SequenceNode {
				for _, entry := range entries.Content {
					orderMapping(entry, order)
				}
			}
		}
	}
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(&root); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func orderMapping(node *yaml.Node, order []string) {
	if node.Kind != yaml.MappingNode {
		return
	}
	ordered := make([]*yaml.Node, 0, len(node.Content))
	for _, key := range order {
		for index := 0; index < len(node.Content); index += 2 {
			if node.Content[index].Value == key {
				ordered = append(ordered, node.Content[index], node.Content[index+1])
			}
		}
	}
	for index := 0; index < len(node.Content); index += 2 {
		if !slices.Contains(order, node.Content[index].Value) {
			ordered = append(ordered, node.Content[index], node.Content[index+1])
		}
	}
	node.Content = ordered
}

func readExisting(path string, noInput bool, confirm func(string) (bool, error)) (map[string]any, os.FileMode, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	if !info.Mode().IsRegular() {
		return nil, 0, fmt.Errorf("%s is not a regular file; refusing to overwrite", path)
	}
	if noInput {
		return nil, 0, fmt.Errorf("%s already exists; refusing to overwrite with --no-input", path)
	}
	confirmed, err := confirm("opslevel.yml already exists. Use its values and overwrite it?")
	if err != nil {
		return nil, 0, err
	}
	if !confirmed {
		return nil, 0, ErrOverwriteDeclined
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}
	var document map[string]any
	if err := yaml.Unmarshal(content, &document); err != nil {
		return nil, 0, fmt.Errorf("could not parse existing opslevel.yml: %w", err)
	}
	if document["version"] != 1 {
		return nil, 0, fmt.Errorf("existing opslevel.yml must use version 1")
	}
	if _, ok := document["component"].(map[string]any); !ok {
		return nil, 0, fmt.Errorf("existing opslevel.yml must contain a component mapping")
	}
	return document, info.Mode().Perm(), nil
}

func applyExistingDefaults(options *Options, fields any) error {
	content, err := yaml.Marshal(fields)
	if err != nil {
		return err
	}
	var previous component
	if err := yaml.Unmarshal(content, &previous); err != nil {
		return fmt.Errorf("could not read existing component values: %w", err)
	}
	for _, field := range []struct {
		value    *string
		previous string
	}{
		{&options.Type, previous.Type}, {&options.Lifecycle, previous.Lifecycle},
		{&options.Tier, previous.Tier}, {&options.Owner, previous.Owner},
		{&options.System, previous.System}, {&options.Language, previous.Language},
		{&options.Framework, previous.Framework}, {&options.Description, previous.Description},
	} {
		if *field.value == "" {
			*field.value = field.previous
		}
	}
	return nil
}

func replaceExisting(path string, content []byte, mode os.FileMode) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".opslevel-*.yml")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(content); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Chmod(file.Name(), mode); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

func completeOptions(
	options *Options,
	askString func(string, string) (string, error),
	askChoice func(string, []string) (string, error),
) error {
	choices := []struct {
		name   string
		value  *string
		values []string
	}{
		{"type", &options.Type, []string{"service", "library", "application", "tool", "configuration", "orchestrator", "test"}},
		{"lifecycle", &options.Lifecycle, []string{"experimental", "development", "production", "deprecated", "retired"}},
		{"tier", &options.Tier, []string{"tier_1", "tier_2", "tier_3", "tier_4"}},
	}
	for _, choice := range choices {
		if !options.NoInput {
			prompt := "Component " + choice.name + ":"
			values := slices.Clone(choice.values)
			if index := slices.Index(values, *choice.value); index >= 0 {
				values = append([]string{*choice.value}, slices.Delete(values, index, index+1)...)
				prompt += " (default: " + *choice.value + ")"
			}
			value, err := askChoice(prompt, values)
			if err != nil {
				return err
			}
			*choice.value = value
		}
		if !slices.Contains(choice.values, *choice.value) {
			return fmt.Errorf("--%s must be one of: %s", choice.name, strings.Join(choice.values, ", "))
		}
	}
	for _, field := range []struct {
		name   string
		prompt string
		value  *string
	}{
		{"owner", "Owning team or chapter alias:", &options.Owner},
		{"system", "System (optional):", &options.System},
		{"language", "Main language (optional):", &options.Language},
		{"framework", "Framework or runtime (optional):", &options.Framework},
		{"description", "Component description:", &options.Description},
		{"package-name", "Package name (required for SonarQube integration):", &options.PackageName},
	} {
		if !options.NoInput {
			var value string
			var err error
			if field.name == "system" {
				value, err = askSystem(*field.value, askString, askChoice)
			} else if field.name == "language" {
				value, err = askLanguage(*field.value, askString, askChoice)
			} else if field.name == "framework" {
				value, err = askFramework(*field.value, askString, askChoice)
			} else {
				value, err = askString(field.prompt, *field.value)
			}
			if err != nil {
				return err
			}
			*field.value = value
		}
		*field.value = strings.TrimSpace(*field.value)
		if *field.value == "" && (field.name == "owner" || field.name == "description" || field.name == "package-name") {
			return fmt.Errorf("--%s is required", field.name)
		}
	}
	return nil
}

func askSystem(
	current string,
	askString func(string, string) (string, error),
	askChoice func(string, []string) (string, error),
) (string, error) {
	values := []string{
		"core_platform",
		"developer_experience",
		"documentation_platform",
		"end_user_portal",
		"identity_and_access_management",
		"market_agent",
		"market_data",
		"market_portal",
		"metering_values",
		"service_platform",
		"settlement",
	}
	return askCustomChoice("System", current, values, askString, askChoice)
}

func askLanguage(
	current string,
	askString func(string, string) (string, error),
	askChoice func(string, []string) (string, error),
) (string, error) {
	return askCustomChoice("Language", current, []string{"Kotlin", "Java", "Python", "Go", "Typescript"}, askString, askChoice)
}

func askFramework(
	current string,
	askString func(string, string) (string, error),
	askChoice func(string, []string) (string, error),
) (string, error) {
	return askCustomChoice("Framework", current, []string{"Ktor", "React", "Ansible"}, askString, askChoice)
}

func askCustomChoice(
	label, current string,
	options []string,
	askString func(string, string) (string, error),
	askChoice func(string, []string) (string, error),
) (string, error) {
	const none = "None"
	const custom = "Custom..."
	values := append([]string{none}, options...)
	values = append(values, custom)
	current = strings.TrimSpace(current)
	defaultValue := none
	if current != "" {
		defaultValue = current
		if index := slices.Index(values, current); index >= 0 {
			values = slices.Delete(values, index, index+1)
		}
		values = append([]string{current}, values...)
	}
	selected, err := askChoice(label+" (default: "+defaultValue+"):", values)
	if err != nil {
		return "", err
	}
	switch selected {
	case none:
		return "", nil
	case custom:
		return askString("Custom "+strings.ToLower(label)+":", current)
	default:
		return selected, nil
	}
}
