# gh-dxp

A GitHub (gh) CLI extension for automating daily development work, brought to you by Elhub's DevXP team. It implements an opinionated workflow based around small and frequent commits, squash merge, and mandatory linting and unit testing.

## Usage

As mentioned, the workflow implemented by this extension is built around the Elhub's development teams opinionated workflow. It is a one-idea, one-commit workflow based around [trunk based development](https://trunkbaseddevelopment.com/).

We **always** lint and unit test commits before submitting them for review in addition to CI/CD checks; there is no point in knowingly pushing bad formatting/bugs into a pull request and wasting the reviewer's time.

To start a change, we create a new local branch:

   ```sh
   gh dxp branch newfeature
   ```

This creates a new branch `newfeature` and switches to it.

We use this new branch to work on our feature. When we feel ready to have our work reviewed, we create a pull request in the working branch.

   ```sh
   gh dxp pr
   ```

This does a number of things. First it runs `lint` and `unit` in order to verify that there are no unneeded issues with the code. If all is well, it will push the `newfeature` branch to GitHub and create a new pull request, asking a few questions as needed to fill in the review template.

When the code is reviewed, run:

   ```sh
   gh dxp merge
   ```

This squash merges your pull request into the main branch and deletes both your local and the remote branch.

In addition, there are some convenience operations available to support our daily work:

   ```sh
   gh dxp lint
   ```

This runs the installed linters on the given repository. Linters need to be configured, either locally in the repository, or with a config file in your home directory.

   ```sh
   gh dxp test
   ```

This runs unit tests on the given repository. The extension will try to guess the unit test framework to run based on project type, but it can also be configured in the configuration files.

To use the status command, simply run:

```sh
gh dxp status
```

This command will display information about the current branch, uncommitted changes, and any changes that are committed but not yet pushed to the remote repository. It's a handy way to ensure that your work is organized and ready for review.

Remember, staying informed about the status of your work helps streamline the development process and facilitates collaboration with your team.

### Configuration options

!!! TODO: Describe configurations options using the .devxp file.

### Aliases

To avoid having to type `gh dxp` constantly, we recommend running:

   ```sh
   gh alias import alias.yml
   ```

On the `alias.yml` file that follows this project. This installs a number of useful aliases for the commands in this extension.

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
  * branch: Create a new local branch, switch to it. If already created, switch to branch
  * pr: Lint, Unit Test and Create a pull request from the existing branch with default info
  * merge: Squash-merge a pull request, switch back to default branch
  * lint: Run linters on project
  * unit: Run unit tests
* Wish-List:
  * Proper Jira integration. Do some basic Jira checks (is ticket assigned, put into progress, etc).

## Linters and Code Analyzers

| Supported             | Language   | Linter       |
|-----------------------|------------|--------------|
| :black_square_button: | Ansible    | ansible-lint |
| :black_square_button: | C#         |              |
| :black_square_button: | CSS        | style-lint   |
| :white_check_mark:    | Golang     | golangci-cli |
| :black_square_button: | Java       | checkstyle   |
| :black_square_button: | Javascript | eslint       |
| :white_check_mark:    | Kotlin     | detekt       |
| :black_square_button: | Markdown   | markdownlint |
| :black_square_button: | OpenAPI    | spectral     |
| :black_square_button: | Shell      | ShellCheck   |
| :black_square_button: | SQL        | sql-lint     |
| :black_square_button: | Terraform  | fmt          |
| :black_square_button: | Typescript | eslint       |
| :white_check_mark:    | YAML       | YamlLint     |
