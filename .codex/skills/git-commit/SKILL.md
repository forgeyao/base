---
name: git-commit
description: Create clean, scoped git commits from an existing worktree. Use when Codex needs to inspect changes, decide what should be staged together, write an imperative commit message, and run the non-destructive git commands needed to produce a reviewable commit.
---

# Git Commit

## Overview

Turn a dirty worktree into one focused commit without guessing. Inspect the current diff, separate unrelated changes, stage only the intended files or hunks, write a concise commit message, and verify the result before finishing.

## Workflow

1. Inspect repository state with non-destructive commands such as `git status --short`, `git diff --stat`, `git diff --cached`, and targeted diffs for changed files.
2. Run unit tests that directly cover the changed code. Do not commit while relevant tests are failing or skipped without explanation.
3. Run `golangci-lint` for the affected Go codebase and fix reported issues before committing.
4. Review the staged diff for secrets or sensitive information before committing, especially tokens, credentials, private endpoints, internal-only data, or machine-specific paths that should not become public.
5. Write the commit message in imperative mood. Keep the subject focused on the user-visible change or engineering intent, not the implementation mechanics alone.
6. Create the commit with a non-interactive git command.
7. Confirm the outcome with `git show --stat --oneline HEAD` or an equivalent summary.

## Staging Rules

- Preserve user changes that are outside the requested commit scope.
- Prefer `git add <path>` when an entire file belongs to the commit.
- Prefer patch-level staging when only part of a file belongs to the commit.
- For newly added hidden files such as `.env.example`, `.github/...`, `.codex/...`, or other dotfiles, confirm they should be committed before staging them.
- Do not rewrite history, amend commits, rebase, or force-push unless the user explicitly asks.
- Do not stage generated files unless they are required outputs of the intended change.

## Commit Message Rules

- Use a short imperative subject line such as `Add JWT service tests` or `Fix Gin response logging`.
- Add a body only when it helps explain motivation, constraints, or follow-up impact.
- Match any repository-specific convention if one is visible in recent history or repo docs.
- Avoid noisy prefixes unless the repository already uses them consistently.

## Required Checks

- Run the unit tests that are directly relevant to the modified code and require them to pass before committing.
- Run `golangci-lint` and fix findings rather than committing known lint failures.
- Inspect the staged patch for sensitive information before commit. Treat this as mandatory for open source repositories.
- Review newly added hidden files separately and confirm they belong in version control before including them in the commit.
- If any required check cannot run, stop and report the blocker instead of pretending the commit is ready.

## Command Pattern

Use commands like these as building blocks and adapt them to the repo:

```bash
git status --short
git diff --stat
git diff --cached
go test ./path/to/changed/package/...
golangci-lint run
git add path/to/file
git add -p
git commit -m "Short imperative subject"
git show --stat --oneline HEAD
```

## Output

Report:

- what was staged
- what validation ran
- whether `golangci-lint` passed
- whether the staged diff was checked for sensitive information
- whether any newly added hidden files were reviewed and confirmed for commit
- the final commit hash and subject
- any remaining unstaged changes that were intentionally left out
