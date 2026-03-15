# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

**cogmit** is a CLI tool that uses AI (via OpenRouter) to generate semantic commit messages from staged git diffs. It follows the Conventional Commits format and prompts the user for confirmation before committing.

## Running

```bash
go run .           # run directly with Go
cogmit             # if installed globally via go install
```

Before running, ensure changes are staged:
```bash
git add <files>
cogmit
```

## Build & Install

```bash
go build -o cogmit.exe .   # build binary locally
go install .               # install to $GOPATH/bin (makes `cogmit` available globally)
```

## Environment

Requires a `.env` file in the project root with:
- `OPENROUTER_API_KEY` — API key for OpenRouter
- `AI_MODEL` — model slug (e.g. `deepseek/deepseek-chat-v3.1`); defaults to `arcee-ai/trinity-large-preview:free`

## Architecture

Go project with the following packages:

- **`main.go`** — entry point; orchestrates the full flow
- **`/ai/openrouter.go`** — sends diff to OpenRouter API and extracts the commit message
- **`/git/git.go`** — wrappers for `git diff --cached`, `git commit`, `git push`, repo root detection
- **`/convention/convention.go`** — loads commit convention from `cogmit-convention.md` in repo root, or uses built-in default
- **`/prompt/prompt.go`** — simple yes/no confirmation prompt via stdin

### Flow

1. Detect repo root via `git rev-parse --show-toplevel`
2. Load convention rules (from `cogmit-convention.md` or built-in default)
3. Get staged diff via `git diff --cached`; exit if nothing staged
4. Call OpenRouter API; prompt requests sections `### Resumo Técnico` and `### Commit Message`
5. Extract commit message (strategy 1: `### Commit Message` section; strategy 2: regex Conventional Commit pattern)
6. Display AI response to user, prompt for confirmation (default: yes)
7. Execute `git commit -m <message>`
8. Prompt to push to remote (default: no)
