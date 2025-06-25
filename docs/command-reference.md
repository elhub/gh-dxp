---
tags: [github, reference]
---

# gh-dxp Command Reference

This document contains documentation for the various commands implemented in `gh-dxp`.

💡 tip:
    You can run any command with the `--help` flag to get more information on how it works. There are lots of
    useful options that are not documented on these pages.


## 🕵️ alias
In order to simplify usage, we have defined a default `alias.yml` that defines the most commonly used workflow
commands. The `alias` command downloads and imports this default file.

**Example:**

```bash
gh dxp alias import
```

By default this "clobbers" (i.e., overwrites) any existing aliases with the same name.


## 🌿 branch
Provides a shortcut for creating and switching to branches.


## ✅ completion
Generates and sets up the autocompletion script for the `gh dxp` command.
This allows you to use tab completion for the various commands and options available in `gh dxp`.
Has options for bash, fish, zsh, and powershell.

## 🆘 help
Provides help for the `gh dxp` command and its subcommands.

## 🧹 lint
Runs MegaLinter on the project. By default, the linter will only run on files that have a diff
to the default branch. If you want to lint everything, you can run the linter using the `--all` flag. Some lint errors
can be fixed by using the `--fix` flag.

### Linter configuration

Configuration of MegaLinter, such as file exclusions or custom rules, is done by adding a `.mega-linter.yml` in the
repository root.

💡 tip:
    Check out [megalinter documentation](https://megalinter.io/7.8.0/configuration/) for more info on how to
    configure the linter.

If no `.mega-linter.yml` is present in the repository root, the lint command will default to using the config defined
in [devxp-lint-configuration](https://github.com/elhub/devxp-lint-configuration). If you want configure the linter,
you should strongly consider whether the changes would be best suited in the
[devxp-lint-configuration](https://github.com/elhub/devxp-lint-configuration) repo instead of your local repository.

In order to keep the default configuration in addition to possible modifications, include the following in your
`mega-linter.yml`:

```yaml
---
EXTENDS:
  - https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml

```

## 🧐 owner
Gets the owner of a specific file or directory. This is useful for determining who to contact if you have questions about the code.

## ⤵️ pr

The `pr` command handles all things related to pull requests.

### pr create

The `pr create` command allows you to create and update diffs/pull requests. By default, it will run both `lint` and
`test` as steps.

**Example:**

```bash
# Start flow to create pr
gh dxp pr create

# Start flow to create pr, but do not run linting and tests
gh dxp pr create --nolint --nounit

# Start flow, with prefilled branch name and commit message
gh dxp pr create -b branchName -m "Add amazing new feature"
```

### pr merge

The `pr merge` command handles the merging of diffs/pull requests.


## 🗃️ repo

Extends the basic repo commands provided by the gh cli.
Also has support for cloning all repositories starting with a given prefix.

### repo clone-all

The `repo clone-all` command allows you to clone some or all repositories within an organization.

**Example:**

```bash
# Clone all repositories from all organization you have access to.
gh dxp repo clone-all

# Clone all repositories that contain "docs"
gh dxp repo clone-all docs

# List the repositories that would be cloned using this command
gh dxp repo clone-all docs --dryrun
```

## 🔎 status
Allows you to get the status of various aspects of the repository, such as existing branches, pull requests, issues etc.

## 📐 template
Generates relevant template files (like .teamcity folder, .gitignore, .editorconfig, etc.) in the current repository.
Also has support for generating base files for gradle projects, if using the `--gradle` flag.

## 🧪 test

The `test` command will attempt to autodetect tests in your current repository and run them. It does so using the following logic:

1. *If* the repository root contains a `Makefile`, the test command will be `make check`
2. *if* the repository root contains a `gradlew`, the test command will be `./gradlew test`
3. *if* the repository root contains a `pom.xml`, the test command will be `mvn test`
4. *if* the repository root contains a `package.json`, the test command will be `npm test`
5. *else* (i.e. none of the above): test will simply print *"no test command found"* and return exit code 0 (success).

💡 tip:
If your setup doesn't neatly fit into any of the options outlined above, you can add a Makefile to your repo and
define the `make check` command however you want.

