package convention

import (
	"os"
	"path/filepath"
)

const defaultConvention = `Use Conventional Commits format: <type>(<scope>): <description>
Types: feat, fix, chore, refactor, style, docs, test, perf, build, ci
Keep the description under 72 characters. Use imperative mood. English only.`

// Load reads cogmit-convention.md from the repo root.
// If the file is not found, it returns the built-in Conventional Commits default.
func Load(repoRoot string) string {
	path := filepath.Join(repoRoot, "cogmit-convention.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConvention
	}
	return string(data)
}
