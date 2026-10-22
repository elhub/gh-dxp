# Copilot instructions

This repository is the `gh-dxp` GitHub CLI extension, written in Go.

## Pull request review conclusions

When reviewing pull requests:

- Apply these wording rules to review summaries, overview comments, headings,
  and individual findings.
- Report findings and material review limitations, not approval or merge
  recommendations.
- Do not state or imply that a pull request is approved or ready to merge.
- Do not use "Approval recommended", "Ready to approve", "Approved", "LGTM",
  "safe to merge", or equivalent verdicts, including in headings.
- When no findings are identified, use exactly "No findings." as the conclusion
  and report any material review limitations separately.
- When findings are identified, use "Findings identified." as the conclusion,
  followed by actionable findings and any material review limitations.
- Keep narrative conclusions separate from GitHub's configured review status.
  These wording rules do not request changes to Copilot approval settings.

## Stack

- Go 1.27.1
- GitHub CLI extension built with Cobra and Bubble Tea
