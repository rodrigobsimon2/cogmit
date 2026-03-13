# Low Level Design — cogmit

## Language & Build

- **Language:** Go (ESM Node.js → Go rewrite)
- **Module:** `cogmit`
- **Binary name:** `cogmit` (`cogmit.exe` on Windows)
- **Build:** `go build -o cogmit .`
- **Distribution:** compiled binary released per platform (Windows/Linux/macOS)

## Package Structure (proposed)

```
cogmit/
├── main.go              # entry point, CLI wiring
├── git/
│   └── git.go           # diff, commit, push via os/exec or go-git
├── ai/
│   └── openrouter.go    # HTTP call to OpenRouter API
├── convention/
│   └── convention.go    # load cogmit-convention.md or return built-in default
└── prompt/
    └── prompt.go        # interactive confirm prompts (commit / push)
```

## Core Functions

### `convention.Load(repoRoot string) → string`
- Looks for `cogmit-convention.md` at `repoRoot/cogmit-convention.md`
- If found: returns file contents as a string
- If not found: returns built-in Conventional Commits default (no warning, no error)

**Built-in default:**
```
Use Conventional Commits format: <type>(<scope>): <description>
Types: feat, fix, chore, refactor, style, docs, test, perf, build, ci
Keep the description under 72 characters. Use imperative mood. English only.
```

### `git.StagedDiff() → (string, error)`
- Runs `git diff --cached`
- Returns error if nothing is staged

### `ai.GenerateCommitMessage(diff, convention string) → (string, error)`
- Calls `POST https://openrouter.ai/api/v1/chat/completions`
- Auth: `Bearer $OPENROUTER_API_KEY`
- Model: `$AI_MODEL` (default: `deepseek/deepseek-chat:free`)
- Temperature: `0.2`, max_tokens: `400`
- Prompt structure:
  ```
  You are a software engineering expert.
  Analyze the following git diff and produce two sections:

  ### Technical Summary
  Explain in up to 3 sentences what changed.

  ### Commit Message
  Write a commit message following these rules:
  <convention content>

  Respond only with these two sections, in English.

  ---
  Diff:
  <diff>
  ```
- Returns the full AI response text

### `ai.ExtractCommitMessage(response string) → (string, error)`
- Strategy 1: find `### Commit Message` section, return first non-header line
- Strategy 2: regex match for Conventional Commit pattern
- Returns error if neither matches

### `prompt.Confirm(question string, defaultYes bool) → bool`
- Prints question to stdout, reads single character input
- Respects default on Enter key

### `main()`
Orchestration:
1. Detect repo root (`git rev-parse --show-toplevel`)
2. `convention.Load(repoRoot)`
3. `git.StagedDiff()` — exit 0 with warning if empty
4. `ai.GenerateCommitMessage(diff, convention)`
5. Print full AI response (summary + commit message)
6. Print commit message prominently (copy-paste friendly)
7. `prompt.Confirm("Commit now?", true)` — exit 0 if no
8. `git.Commit(commitMessage)`
9. `prompt.Confirm("Push to remote?", false)` — skip if no
10. `git.Push()` — exit 1 on failure

## Environment Variables

Loaded via `os.Getenv` — no `.env` file. Users set these in their shell profile (`~/.bashrc`, `~/.zshrc`, Windows system env, etc.).

| Variable | Required | Default | Description |
|---|---|---|---|
| `OPENROUTER_API_KEY` | Yes | — | OpenRouter API key |
| `AI_MODEL` | No | `deepseek/deepseek-chat:free` | LLM model slug |

## Error Handling

| Condition | Behavior |
|---|---|
| No staged changes | Exit 0 with yellow warning |
| `cogmit-convention.md` not found | Silently use built-in default |
| AI response malformed | Exit 1 with error details |
| Commit message not extractable | Exit 1 with hint |
| User cancels commit | Exit 0 |
| User cancels push | Exit 0 (commit already done) |
| Push fails | Exit 1 with error |
