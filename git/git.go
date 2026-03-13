package git

import (
	"errors"
	"os/exec"
	"strings"
)

// RepoRoot returns the absolute path to the root of the current git repository.
func RepoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", errors.New("not inside a git repository")
	}
	return strings.TrimSpace(string(out)), nil
}

// StagedDiff returns the output of `git diff --cached`.
// Returns an error if there are no staged changes.
func StagedDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		return "", err
	}
	diff := strings.TrimSpace(string(out))
	if diff == "" {
		return "", errors.New("no staged changes — run `git add` first")
	}
	return diff, nil
}

// Commit creates a git commit with the given message.
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = nil
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}

// Push pushes the current branch to its tracking remote.
func Push() error {
	out, err := exec.Command("git", "push").CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
