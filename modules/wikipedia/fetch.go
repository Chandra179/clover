package wikipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Chandra179/gosdk/logger"
)

var CategorySearchTerms = map[string]string{
	"economy": "economy OR economics OR finance OR business",
	"tech":    "technology OR computing OR software OR internet",
}

type CategoryResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
}

func (d *Dependencies) FetchCategory(category string) ([]CategoryResult, error) {
	ctx := context.Background()

	query, ok := CategorySearchTerms[category]
	if !ok {
		return nil, fmt.Errorf("unknown category: %s (allowed: economy, tech)", category)
	}

	apiURL := fmt.Sprintf("%s/w/api.php", d.Config.Wikipedia.BaseURL)
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srlimit", "10")
	params.Set("format", "json")

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("wikipedia: build request: %w", err)
	}
	req.Header.Set("User-Agent", d.Config.Wikipedia.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wikipedia: request failed: %w", err)
	}
	defer resp.Body.Close()

	var qr QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&qr); err != nil {
		return nil, fmt.Errorf("wikipedia: decode failed: %w", err)
	}
	if qr.Error != nil {
		return nil, fmt.Errorf("wikipedia: api error: %s", qr.Error.Info)
	}
	if qr.Query == nil {
		return nil, fmt.Errorf("wikipedia: empty query result")
	}

	var results []CategoryResult
	for _, s := range qr.Query.Search {
		u := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(s.Title))
		results = append(results, CategoryResult{
			Title:       s.Title,
			URL:         u,
			Content:     stripHTML(s.Snippet),
			Category:    category,
			Source:      "wikipedia",
			PublishedAt: s.Timestamp,
		})
	}

	d.Logger.Info(ctx, "wikipedia: fetch done",
		logger.Field{Key: "category", Value: category},
		logger.Field{Key: "count", Value: len(results)},
	)

	return results, nil
}

func stripHTML(s string) string {
	var result []byte
	inTag := false
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			inTag = true
			continue
		}
		if s[i] == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result = append(result, s[i])
		}
	}
	return string(result)
}
