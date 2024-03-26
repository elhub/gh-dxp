# gh-dxp

A gh client extension brought to you by Elhub's DevXP team.

## Thoughts

* Have a .devxp file (YAML) in repo that we can read config details from
  * Linters to use
* Nice to have:
  * new:  Create a new project with default template files.
  * work: Create a new local branch, switch to it. If already created, switch to branch
  * diff: Lint, Unit Test and Create a pull request from the existing branch with default info
  * land: Squash-merge a pull request, switch back to default branch
  * lint: Run linters on project
  * unit: Run unit tests
* Dreams:
  * Proper Jira integration. Do some basic Jira checks (is ticket assigned, put into progress, etc).

## Linters and Code Analyzers

| Supported          | Language    | Linter          |
| ------------------ | ----------- | --------------- |
| :white_check_mark: | Ansible     | ansible-lint    |
| :white_check_mark: | C#          |                 |
| :white_check_mark: | CSS         | style-lint      |
| :heavy_check_mark: | Golang      | golangci-cli    |
| :white_check_mark: | Java        | checkstyle      |
| :white_check_mark: | Javascript  | eslint          |
| :heavy_check_mark: | Kotlin      | detekt          |
| :white_check_mark: | Markdown    | markdownlint    |
| :white_check_mark: | OpenAPI     | spectral        |
| :white_check_mark: | Shell       | ShellCheck      |
| :white_check_mark: | SQL         | sql-lint        |
| :white_check_mark: | Terraform   | fmt             |
| :white_check_mark: | Typescript  | eslint          |
| :heavy_check_mark: | YAML        | YamlLint        |
