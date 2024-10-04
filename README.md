# gh-dxp

A GitHub (gh) CLI extension for automating daily development work, brought to you by Elhub's DevXP team. It implements an opinionated workflow based around small and frequent commits, squash merge, and mandatory linting and unit testing. To view more detailed documentation, please refer to the gh-dxp page in docs-support.

### User Guide
Using the `-h` flag with any command will display relevant documentation.
In addition to that, a user guide for `gh dxp` can be found [here](https://docs.elhub.cloud/support/applications/gh-dxp/index.html)

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
   <summary><strong>Install from source</strong></summary>

If you want to install this extension **from source**, such as for development, follow these steps:

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
