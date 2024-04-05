# gh-dxp

A GitHub (gh) CLI extension for automating daily development work, brought to you by Elhub's DevXP team. It implements an opinionated workflow based around small and frequent commits, squash-merge, and mandatory linting and unit testing.

## Usage







## Installation

1. [Install the `gh` CLI](https://github.com/cli/cli#installation)
2. Install gh-dxp:
   ```sh
   gh extension install elhub/gh-dxp
   ```

<details>
   <summary><strong>Manual Install</strong></summary>

If you want to install this extension **manually**, follow these steps:

1. Clone the repo

   ```bash
   # git
   git clone https://github.com/elhub/gh-dxp
   ```

2. Build and install locally

   ```bash
   cd gh-dxp; make clean install
   ```
</details>

## RoadMap

Following are some of the things we are thinking/working on.

* Settings files that can be used to configure linters, etc.
  * Linters to use
* Workflows:
  * new:  Create a new project with default template files.
  * work: Create a new local branch, switch to it. If already created, switch to branch
  * diff: Lint, Unit Test and Create a pull request from the existing branch with default info
  * land: Squash-merge a pull request, switch back to default branch
  * lint: Run linters on project
  * unit: Run unit tests
* Would-Love-To-Have:
  * Proper Jira integration. Do some basic Jira checks (is ticket assigned, put into progress, etc).

## Linters and Code Analyzers

| Supported          | Language    | Linter          |
| ------------------ | ----------- | --------------- |
| :white_square: | Ansible     | ansible-lint    |
| :white_square: | C#          |                 |
| :white_square: | CSS         | style-lint      |
| :white_check_mark: | Golang      | golangci-cli    |
| :white_square: | Java        | checkstyle      |
| :white_square: | Javascript  | eslint          |
| :white_check_mark: | Kotlin      | detekt          |
| :white_square: | Markdown    | markdownlint    |
| :white_square: | OpenAPI     | spectral        |
| :white_square: | Shell       | ShellCheck      |
| :white_square: | SQL         | sql-lint        |
| :white_square: | Terraform   | fmt             |
| :white_square: | Typescript  | eslint          |
| :white_check_mark: | YAML        | YamlLint        |
