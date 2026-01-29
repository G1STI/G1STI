package app

import (
	"fmt"
	"os"
	"strings"
)

func readLogTail(path string, lines int) (string, error) {
	if lines <= 0 {
		lines = 200
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read log: %w", err)
	}
	content := strings.ReplaceAll(string(payload), "\r\n", "\n")
	parts := strings.Split(content, "\n")
	if len(parts) <= lines {
		return strings.TrimSpace(content), nil
	}
	start := len(parts) - lines
	return strings.TrimSpace(strings.Join(parts[start:], "\n")), nil
}
