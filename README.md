# gh-dxp

A GitHub (gh) CLI extension for automating daily development work, brought to you by Elhub's DevXP team. It implements an opinionated workflow based around
small and frequent commits, squash merge, and mandatory linting and unit testing. To view more detailed documentation, please refer to the gh-dxp page
in docs-support.

## User Guide
Using the `-h` flag with any command will display relevant documentation.
In addition to that, a [user guide](https://docs.elhub.cloud/support/applications/gh-dxp/index.html) for `gh dxp` is also available.

## Jira Integration (Optional)

`gh dxp pr create` can suggest related Jira tickets based on your commits, PR title and description.
The integration is opt-in and uses your own Jira credentials.

### Setup

1. Create an API token at [id.atlassian.com/manage-api-tokens](https://id.atlassian.com/manage-api-tokens).
2. Provide your credentials using **one** of the following (environment variables take priority):

   | Method | Configuration |
   |--------|---------------|
   | File   | `~/.jira_token` containing a single line: `email:token` |
   | Env    | `JIRA_USERNAME` and `JIRA_API_TOKEN` |

   ```sh
   echo "your-email@example.com:your_api_token" > ~/.jira_token
   chmod 600 ~/.jira_token
   ```

To opt out and never be asked about Jira, run `echo disabled > ~/.jira_token`.

### Behavior

- Runs only when creating a **new** PR without `--issues`. Updating an existing PR does not query Jira.
- Searches your open, assigned tickets (the `ET` project is excluded) and lists the best matches.
- Reads only ticket key, summary and description. Nothing else is accessed or stored.
- Without credentials, the step is skipped and PR creation works as usual.
- Requests time out after 5 seconds, so Jira issues never block PR creation.

### Troubleshooting

- Verify credentials: `curl -u "email:token" https://elhub.atlassian.net/rest/api/3/myself`
- Make sure `~/.jira_token` contains a single `email:token` line.
- No suggestions usually means none of your open tickets match the commit text.
## Aliases

To avoid having to type `gh dxp` constantly, we recommend running:

   ```sh
   gh alias import alias.yml
   ```

The `alias.yml` file included in this project installs a number of useful aliases for the commands in this extension.

## Installation

1. [Install the `gh` CLI](https://github.com/cli/cli#installation)
2. Install gh-dxp:
    ```sh
    gh extension install elhub/gh-dxp
    ```

<details>
   <summary><strong>Install from source</strong></summary>

If you want to install this extension **from source** for development, follow these two steps:

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
