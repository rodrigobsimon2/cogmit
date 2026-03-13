package main

import (
	"cogmit/ai"
	"cogmit/convention"
	gogit "cogmit/git"
	"cogmit/prompt"
	"fmt"
	"os"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

func main() {
	// 1. Detect repo root
	repoRoot, err := gogit.RepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, colorRed+"✗ "+err.Error()+colorReset)
		os.Exit(1)
	}

	// 2. Load convention rules
	conv := convention.Load(repoRoot)

	// 3. Get staged diff
	diff, err := gogit.StagedDiff()
	if err != nil {
		fmt.Println(colorYellow + "⚠  " + err.Error() + colorReset)
		os.Exit(0)
	}

	// 4. Generate commit message via AI
	fmt.Println(colorCyan + "\n🤖 Generating commit suggestion..." + colorReset)
	aiResponse, err := ai.GenerateCommitMessage(diff, conv)
	if err != nil {
		fmt.Fprintln(os.Stderr, colorRed+"✗ AI error: "+err.Error()+colorReset)
		os.Exit(1)
	}

	// 5. Display full AI response
	fmt.Println("\n" + aiResponse)

	// 6. Extract and highlight the commit message
	commitMsg, err := ai.ExtractCommitMessage(aiResponse)
	if err != nil {
		fmt.Fprintln(os.Stderr, colorRed+"✗ "+err.Error()+colorReset)
		os.Exit(1)
	}

	fmt.Println("\n" + colorBold + "┌─ Commit message " + strings.Repeat("─", 40) + colorReset)
	fmt.Println(colorBold + "│  " + commitMsg + colorReset)
	fmt.Println(colorBold + "└" + strings.Repeat("─", 57) + colorReset)

	// 7. Confirm commit
	if !prompt.Confirm("\nCommit now?", true) {
		fmt.Println(colorYellow + "✗ Commit cancelled." + colorReset)
		os.Exit(0)
	}

	// 8. Commit
	if err := gogit.Commit(commitMsg); err != nil {
		fmt.Fprintln(os.Stderr, colorRed+"✗ Commit failed: "+err.Error()+colorReset)
		os.Exit(1)
	}
	fmt.Println(colorGreen + "✔ Committed: " + commitMsg + colorReset)

	// 9. Confirm push
	if !prompt.Confirm("Push to remote?", false) {
		os.Exit(0)
	}

	// 10. Push
	if err := gogit.Push(); err != nil {
		fmt.Fprintln(os.Stderr, colorRed+"✗ Push failed: "+err.Error()+colorReset)
		os.Exit(1)
	}
	fmt.Println(colorGreen + "✔ Pushed." + colorReset)
}
