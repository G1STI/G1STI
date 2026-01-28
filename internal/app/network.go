package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ipResponse struct {
	IP string `json:"ip"`
}

func fetchPublicIP(ctx context.Context) (string, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org?format=json", nil)
	if err != nil {
		return "", fmt.Errorf("create ip request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch ip: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ip response status %d", resp.StatusCode)
	}

	var payload ipResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode ip response: %w", err)
	}
	if payload.IP == "" {
		return "", fmt.Errorf("ip response empty")
	}
	return payload.IP, nil
}
