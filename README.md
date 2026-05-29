# gh-dxp

`gh-dxp` is a GitHub CLI (`gh`) extension from Elhub DevXP for automating common day-to-day repository workflows.

The extension is opinionated toward:
1. Small and frequent commits
2. Mandatory linting and tests before PRs
3. Squash-merge focused PR flow

## Installation

1. [Install the `gh` CLI](https://github.com/cli/cli#installation)
2. Install the extension:

```sh
gh extension install elhub/gh-dxp
```

3. Verify installation:

```sh
gh dxp --help
```

<details>
<summary><strong>Install from source (development)</strong></summary>

```sh
git clone https://github.com/elhub/gh-dxp
cd gh-dxp
make clean install
```

</details>

## Quick start

Import recommended aliases from this repository:

```sh
gh alias import alias.yml
```

Common first commands:

```sh
gh dxp status
gh dxp test
gh dxp lint
gh dxp pr create
```

## Command overview

`gh dxp` includes commands for:
- `alias`: import and manage extension aliases
- `branch`: create and switch branches quickly
- `completion`: shell completion setup
- `lint`: run MegaLinter (supports `--all` and `--fix`)
- `owner`: identify owner for a file or directory
- `pr`: create/update and merge pull requests
- `repo`: repository utilities (including `clone-all`)
- `status`: repo status overview
- `template`: generate common repo templates
- `test`: run tests with repository auto-detection

For complete command docs, see:
- `gh dxp <command> --help`
- [`docs/command-reference.md`](docs/command-reference.md)
- [User guide (docs-support)](https://docs.elhub.cloud/support/applications/gh-dxp/index.html)

## Behavior notes

- `gh dxp test` auto-detects how to run tests in this order:
  1. `make check`
  2. `./gradlew test`
  3. `mvn test`
  4. `npm test`
- `gh dxp lint` uses local `.mega-linter.yml` when present; otherwise it falls back to
-  [devxp-lint-configuration](https://github.com/elhub/devxp-lint-configuration).

## Development

Local development workflow:

```sh
make dep
make check
make vet
make build
```

Install your local build as an extension:

```sh
make install
```
