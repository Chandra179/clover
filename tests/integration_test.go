//go:build integration

package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

type newsResponse struct {
	Category string         `json:"category"`
	Count    int            `json:"count"`
	Sources  map[string]int `json:"sources"`
	Items    []struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		Source      string `json:"source"`
		PublishedAt string `json:"published_at"`
	} `json:"items"`
}

func TestIntegration_NewsEndpoint(t *testing.T) {
	base := "http://localhost:8001"
	client := &http.Client{Timeout: 30 * time.Second}

	tests := []struct {
		name     string
		path     string
		wantCode int
	}{
		{"default category", "/news", http.StatusOK},
		{"tech category", "/news?category=tech", http.StatusOK},
		{"economy category", "/news?category=economy", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(base + tt.path)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Fatalf("expected %d, got %d", tt.wantCode, resp.StatusCode)
			}

			var nr newsResponse
			if err := json.NewDecoder(resp.Body).Decode(&nr); err != nil {
				t.Fatalf("decode: %v", err)
			}

			if nr.Category == "" {
				t.Error("empty category in response")
			}
			if nr.Count < 0 {
				t.Error("negative count")
			}
		})
	}
}
