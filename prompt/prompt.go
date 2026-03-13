package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm prints a yes/no question and reads the user's answer.
// If defaultYes is true, Enter defaults to yes (shown as Y/n).
// If defaultYes is false, Enter defaults to no (shown as y/N).
func Confirm(question string, defaultYes bool) bool {
	hint := "y/N"
	if defaultYes {
		hint = "Y/n"
	}
	fmt.Printf("%s (%s): ", question, hint)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}
