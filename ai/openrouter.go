package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	apiURL       = "https://openrouter.ai/api/v1/chat/completions"
	defaultModel = "deepseek/deepseek-chat:free"
)

type requestBody struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseBody struct {
	Choices []struct {
		Message struct {
			Content   string `json:"content"`
			Reasoning string `json:"reasoning"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message  string `json:"message"`
		Metadata *struct {
			Raw string `json:"raw"`
		} `json:"metadata"`
	} `json:"error"`
}

// GenerateCommitMessage sends the staged diff and convention rules to OpenRouter
// and returns the full AI response text.
func GenerateCommitMessage(diff, convention string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENROUTER_API_KEY environment variable is not set")
	}

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = defaultModel
	}

	prompt := fmt.Sprintf(`You are a software engineering expert.
Analyze the following git diff and produce two sections:

### Technical Summary
Explain in up to 3 sentences what changed.

### Commit Message
Write a commit message following these rules:
%s

Respond only with these two sections, in English.

---
Diff:
%s`, convention, diff)

	body := requestBody{
		Model:       model,
		Messages:    []message{{Role: "user", Content: prompt}},
		Temperature: 0.2,
		MaxTokens:   400,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result responseBody
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	if result.Error != nil {
		msg := result.Error.Message
		if result.Error.Metadata != nil && result.Error.Metadata.Raw != "" {
			msg = result.Error.Metadata.Raw
		}
		return "", fmt.Errorf("AI API error: %s", msg)
	}

	if len(result.Choices) == 0 {
		return "", errors.New("AI returned an empty response")
	}

	// Some reasoning models return empty content with the response in reasoning.
	content := strings.TrimSpace(result.Choices[0].Message.Content)
	if content == "" {
		content = strings.TrimSpace(result.Choices[0].Message.Reasoning)
	}
	if content == "" {
		return "", errors.New("AI returned an empty response")
	}

	return content, nil
}

var conventionalCommitRe = regexp.MustCompile(
	`(?i)(feat|fix|chore|refactor|style|docs|test|perf|build|ci)(\([^)]+\))?:\s.+`,
)

// ExtractCommitMessage parses the commit message from the AI response.
// Strategy 1: looks for the "### Commit Message" section.
// Strategy 2: falls back to a Conventional Commit regex match.
func ExtractCommitMessage(response string) (string, error) {
	// Strategy 1: section header
	if idx := strings.Index(strings.ToLower(response), "### commit message"); idx != -1 {
		rest := response[idx:]
		for _, line := range strings.Split(rest, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			return line, nil
		}
	}

	// Strategy 2: regex
	if match := conventionalCommitRe.FindString(response); match != "" {
		return strings.TrimSpace(match), nil
	}

	return "", errors.New("could not extract a commit message — check that the AI returned a Conventional Commits format")
}
