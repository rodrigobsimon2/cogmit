# High Level Design — cogmit

## Purpose

cogmit is a Go CLI tool that enforces commit message conventions across teams by automating generation with AI. It reads staged changes and an optional per-repo convention file (`cogmit-convention.md`), then uses an LLM to produce an English commit message the developer can copy-paste or commit directly.

## Goals

- Enforce commit message format defined by the repo's `cogmit-convention.md`
- Fall back to standard Conventional Commits when no convention file is present
- Output English commit messages only
- Distribute as a single static binary (no runtime dependencies)
- Keep the developer in control: display message first, then offer to commit

## Non-Goals

- Does not validate or lint existing commit history
- Does not manage branches, PRs, or changelogs
- Does not support multi-language output (English only)
- Distribution strategy not yet decided

## System Overview

```
Developer
   │
   ├─ git add <files>
   └─ cogmit
         │
         ├─ read cogmit-convention.md (from repo root, optional)
         ├─ git diff --cached                  (reads staged diff)
         ├─ OpenRouter API (LLM)               (generates summary + commit message in English)
         ├─ display message                    (user can copy-paste)
         ├─ prompt: "Commit now? (y/N)"
         │     ├─ yes → simple-git commit
         │     └─ no  → exit 0
         └─ prompt: "Push to remote? (y/N)"   (spec-002)
               ├─ yes → git push
               └─ no  → exit 0
```

## Key Design Decisions

- **Go** — single static binary, easy to distribute via `go install` or GitHub release artifacts; no Node/npm required in target environments
- **Per-repo convention file** — `cogmit-convention.md` in the repo root customizes the AI prompt; teams commit this file to enforce their own rules
- **Graceful fallback** — if `cogmit-convention.md` is absent, built-in Conventional Commits defaults are used silently
- **Display-first UX** — message is always shown before any git action so the developer can copy it regardless of choice
- **OpenRouter as AI gateway** — model is swappable via env var without code changes

## Future Directions

- `cogmit --push` flag or config option to skip the push prompt
- Shareable per-repo config file (`.cogmitrc`) for model, tone, scope prefixes
- Inline message editing before committing
