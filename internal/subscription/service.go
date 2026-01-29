package subscription

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"g1sti/internal/model"
	"g1sti/internal/parser"
)

type Service struct {
	client *http.Client
}

func NewService(client *http.Client) *Service {
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	return &Service{client: client}
}

func (s *Service) Fetch(ctx context.Context, subscriptionURL string, subscriptionID string) ([]model.Node, string, error) {
	if subscriptionURL == "" {
		return nil, "", fmt.Errorf("subscription url is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, subscriptionURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create subscription request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fetch subscription: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("subscription response status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read subscription body: %w", err)
	}

	checksum := sha256.Sum256(body)
	checksumHex := hex.EncodeToString(checksum[:])

	lines := decodeSubscriptionLines(string(body))
	nodes := make([]model.Node, 0, len(lines))
	for _, line := range lines {
		parsed, err := parser.ParseURI(line)
		if err != nil {
			continue
		}
		nodes = append(nodes, model.Node{
			ID:             newID(),
			SubscriptionID: subscriptionID,
			Tag:            parsed.Tag,
			Type:           parsed.Type,
			Server:         parsed.Server,
			Port:           parsed.Port,
			Raw:            parsed.Raw,
			ParsedFields:   parsed.ParsedFields,
		})
	}

	return nodes, checksumHex, nil
}

func decodeSubscriptionLines(payload string) []string {
	trimmed := strings.TrimSpace(payload)
	if trimmed == "" {
		return nil
	}

	if strings.Contains(trimmed, "://") {
		return splitLines(trimmed)
	}

	decoded, err := parser.DecodeBase64String(trimmed)
	if err == nil && strings.Contains(decoded, "://") {
		return splitLines(decoded)
	}

	return splitLines(trimmed)
}

func splitLines(payload string) []string {
	lines := strings.Split(payload, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		clean := strings.TrimSpace(line)
		if clean == "" {
			continue
		}
		result = append(result, clean)
	}
	return result
}

func newID() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
