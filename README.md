# cogmit

CLI tool that uses AI to generate [Conventional Commits](https://www.conventionalcommits.org/) messages from your staged git diffs. Powered by [OpenRouter](https://openrouter.ai/).

## Install

Download a binary from the [releases page](https://github.com/rodrigobsimon2/cogmit/releases), or install with Go:

```bash
go install github.com/rodrigobsimon2/cogmit@latest
```

## Setup

Set your OpenRouter API key via a `.env` file in the repo root:

```
OPENROUTER_API_KEY=your-key-here
AI_MODEL=deepseek/deepseek-chat-v3.1   # optional
```

Defaults to `arcee-ai/trinity-large-preview:free` if `AI_MODEL` is not set.

## Usage

```bash
git add <files>
cogmit
```

cogmit reads your staged diff, sends it to the AI model, and proposes a commit message. You confirm before it commits, and optionally pushes to remote.

## Custom conventions

Add a `cogmit-convention.md` file to your repo root to customize commit message rules. Otherwise a built-in Conventional Commits default is used.

## License

MIT
