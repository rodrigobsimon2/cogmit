# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

**cogmit** is a CLI tool that uses AI (via OpenRouter) to generate semantic commit messages from staged git diffs. It follows the Conventional Commits format and prompts the user for confirmation before committing.

## Running

```bash
npm start          # run directly with node
node index.js      # equivalent
cogmit             # if installed globally via npm link
```

Before running, ensure changes are staged:
```bash
git add <files>
node index.js
```

## Environment

Requires a `.env` file in the project root with:
- `OPENROUTER_API_KEY` — API key for OpenRouter
- `AI_MODEL` — model slug (e.g. `deepseek/deepseek-chat-v3.1`)

## Architecture

Single-file ESM script (`index.js`) with this flow:

1. **`getGitDiff()`** — runs `git diff --cached`, exits if nothing is staged
2. **`getAiSummary(diff)`** — sends diff to OpenRouter API; the prompt requests a Markdown response with two sections: `### Resumo Técnico` and `### Mensagem de Commit`
3. **`extractCommitMessage(aiResponse)`** — parses the AI response, first by looking for the markdown section, then by regex-matching a Conventional Commit pattern (`feat|fix|chore|...`)
4. **`main()`** — orchestrates the above, shows the full AI response to the user, then uses `inquirer` to confirm before calling `simple-git` to execute the commit

## Key Dependencies

- `simple-git` — executes the final `git.commit()`
- `inquirer` — interactive confirmation prompt
- `chalk` — colored terminal output
- `dotenv` — loads `.env` from the script's own directory (not `process.cwd()`)
