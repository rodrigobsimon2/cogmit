# Spec 001 — Foundation (Go Rewrite)

**Status:** Implemented
**Created:** 2026-03-12

## Context

cogmit was originally written in Node.js. It is being rewritten in Go to produce a single static binary suitable for developer tooling and DevOps environments — no runtime required, easy to install via `go install` or a GitHub release.

## Problem Statement

The Node.js implementation requires npm/Node installed on the developer's machine and cannot be distributed as a standalone binary. For a CLI tool intended for team-wide adoption across diverse environments, Go is a better fit.

Additionally, teams need a way to define their own commit message conventions per repo, and all output should be in English.

## Proposed Solution

Rewrite `index.js` in Go with the following behavior:

1. Detect the repo root via `git rev-parse --show-toplevel`
2. Load `cogmit-convention.md` from the repo root (fall back to built-in Conventional Commits default if absent)
3. Run `git diff --cached` to get staged changes
4. Send the diff + convention rules to OpenRouter API, requesting an English response with two sections: `### Technical Summary` and `### Commit Message`
5. Display the full AI response (summary + commit message) — user can copy the message
6. Prompt: `Commit now? (Y/n)` — commit on yes, exit 0 on no
7. Prompt: `Push to remote? (y/N)` — push on yes, exit 0 on no

## Acceptance Criteria

- [x] Builds to a single static binary (`go build -o cogmit .`)
- [x] Reads `cogmit-convention.md` from `git rev-parse --show-toplevel` if present
- [x] Falls back to built-in Conventional Commits convention silently if file is absent
- [x] Calls OpenRouter API with diff + convention in English prompt
- [x] AI response is always in English
- [x] Displays technical summary and commit message to stdout
- [x] Prompts user to confirm commit (default: yes)
- [x] Prompts user to confirm push (default: no)
- [x] Exits 0 on cancellation at any step
- [x] Exits 1 with clear message on AI or git errors

## Package Structure

```
main.go
git/git.go
ai/openrouter.go
convention/convention.go
prompt/prompt.go
```

## Out of Scope

- Config file (`.cogmitrc`) support
- Multi-language output
- `--push` / `--no-push` flags
- Inline message editing
