# PR Bot Workflow Guide

This document explains how the pull request automation works with GitHub Actions to validate and merge PRs using maintainer commands.

## Commands

The bot responds to two main commands in PR comments:

1. **`/ok-to-test`** - Runs validation (tests and linting) on a PR
2. **`/approve`** - Approves and merges a validated PR

## Workflow Steps

### Validation Process (`/ok-to-test`)

1. A maintainer comments `/ok-to-test` on a PR
2. The bot verifies the commenter is a maintainer
3. The bot runs tests and linting on the PR code
4. The bot posts detailed results as a comment with test and lint output
5. The bot adds a label based on results:
   - `ready-to-merge` if all checks pass
   - `needs-work` if any checks fail

### Approval Process (`/approve`)

1. A maintainer comments `/approve` on a PR
2. The bot verifies the commenter is a maintainer
3. The bot checks if the PR has the `ready-to-merge` label
4. If validated, the bot attempts to merge the PR (squash merge)
5. The bot posts a comment about the merge result (success or failure)

## Maintainer Configuration

Maintainers are defined in `.github/maintainers.yml`:

```yaml
maintainers:
  - username1
  - username2
  # Add GitHub usernames of maintainers here
```

Only users listed in this file can trigger the validation and approval commands.

## Labels

The workflow uses these labels to track PR status:

- **`ready-to-merge`**: PR has passed validation
- **`needs-work`**: PR has failed validation

## Example Flow

1. Contributor submits a PR
2. Maintainer comments `/ok-to-test`
3. Bot runs tests and linting
4. Bot posts results and adds appropriate label
5. If validation passes, maintainer comments `/approve`
6. Bot merges the PR

## Troubleshooting

- If validation fails, fix the issues and have a maintainer comment `/ok-to-test` again
- If merge fails due to conflicts, resolve them and try `/approve` again
- If the bot doesn't respond to a command, ensure:
  - The commenter is listed in maintainers.yml
  - GitHub Actions workflows are enabled for the repository
  - The GitHub token has sufficient permissions

## Required Permissions

For the workflows to function, ensure your repository has:

- GitHub Actions enabled
- Workflow permissions set to allow write access to pull requests and contents