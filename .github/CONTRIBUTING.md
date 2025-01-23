# Contributing

So you want to contribute ? Awesome. ❤️

All types of contributions are encouraged and valued.  See below for different ways to help and details about how
this project handles them. Please make sure to read the relevant section before making your contribution. It will
make it a lot easier for us maintainers and smooth out the experience for everyone involved. We look forward to
your contributions. 🎉

> If you like the project, but just don't have time to contribute, that's fine. There are other easy ways to support
> the project and show your appreciation, which we would also be very happy about:
>
> * Star the project
> * Refer this project in your project's readme
> * Mention the project on social media, meetups, etc.

## Code of Conduct

We adhere to the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).
By participating, you are expected to uphold this code.

## Issue Tracker

Elhub employees should generally use our internal Jira instance (follow the link in the [README](../README.md));
everyone else is welcome to [open an issue on GitHub](/issues/).

When writing issues:

* Ensure you are running on the latest version.
* Always check whether there is a previous, similar issue. If there is, consider commenting on that issue rather
  than opening an entirely new issue.
* Provide as much context as you can about what you're running into. Describe what you're observing and what you
  expected, steps to reproduce the issue, technical details, etc.

> Do not, **ever**, report security related issues, vulnerabilities or bugs including sensitive information to the
> public issue tracker or elsewhere in public. For sensitive bugs, <email:security@elhub.no>.

## I Have a Question

> If you want to ask a question, we assume that you have read the available documentation (start with ../README.md).

Post your question as a new issue or comment on an existing suitable issue.

We will try to respond as soon as possible.

## I Want To Contribute

> ### Legal Notice
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the
> necessary rights to the content and that the content you contribute may be provided under the project license.

### Reporting Bugs

Before submitting a bug report, please investigate as much as you can beforehand, focusing especially on isolating
the problem and understanding how to recreate it. Then submit the bug in the [issue tracker](#issue-tracker).

### Suggesting Enhancements

Before submitting an enhancement, consider whether the idea fits within the scope and aims of the project.
[Asking a question](#i-have-a-question) can be a good starting point. Then submit the enchancement as a suggestion
through the [issue tracker](#issue-tracker).

### Your First Code Contribution

Ideally, use our GitHub CLI extension ([gh-dxp](https://github.com/elhub/gh-dxp)) to enforce linting/style rules
and format your pull request. Otherwise, use our
[pull request template](https://github.com/elhub/devxp-project-template/blob/main/resources/.github/pull_request_template.md)
as a starting point.

The title should use present tense ("Add feature" not "Added feature") and use the imperative mood.

All our code is linted with [MegaLinter](https://megalinter.io). To run linting manually, use:

```bash
npx mega-linter-runner --install
npx mega-linter-runner --flavor cupcake -e MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml
```

To run the local tests, check the [README](../README.md).

Providing linted _and_ tested code significantly increases the chance that a pull request will be accepted.

### Improving The Documentation

Markdown is our format of preference for all documentation. Submitting new or corrections to documentation follow the
same procedure as for code. You do not need to run unit tests on documentation, but you should run the linter (which
includes [markdownlint](https://github.com/DavidAnson/markdownlint)).
