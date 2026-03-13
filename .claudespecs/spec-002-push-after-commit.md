# Spec 002 — Push After Commit

**Status:** Draft
**Created:** 2026-03-12

## Context

After cogmit successfully creates a commit, the natural next step is often `git push`. Since cogmit already has control of the terminal and git context at that point, it can offer the push as a follow-up action within the same flow.

## Problem Statement

Running `git push` manually after cogmit finishes breaks the flow. The developer must switch focus to the terminal, remember the remote/branch, and run the command separately.

## Proposed Solution

After a successful commit, prompt the user:

```
Push to remote? (y/N)
```

- Default is **No** — push is opt-in per run
- If confirmed, run `git push` (current branch, default remote)
- If push fails, surface the error and exit 1
- If the user declines, exit 0 (commit is already done — no rollback)

## Acceptance Criteria

- [ ] After a successful commit, push confirmation prompt appears
- [ ] Default answer is No (pressing Enter skips push)
- [ ] On confirmation, `git push` executes against the current branch's tracking remote
- [ ] On push success, green success message is shown
- [ ] On push failure, error is printed and process exits 1
- [ ] On user cancellation, process exits 0 cleanly

## Out of Scope

- Selecting a specific remote or branch
- Force push support
- `--push` / `--no-push` CLI flags (future spec)
