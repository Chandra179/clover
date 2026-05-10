package hackernews

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

var CategoryPrefixes = map[string]string{
	"economy": "",
	"tech":    "",
}

func (d *Dependencies) FetchCategory(category string) ([]CategoryResult, error) {
	ctx := context.Background()

	query, ok := CategorySearchTerms[category]
	if !ok {
		return nil, fmt.Errorf("hackernews: unknown category: %s (allowed: economy, tech)", category)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	apiURL := fmt.Sprintf("https://hn.algolia.com/api/v1/search?query=%s&tags=story&hitsPerPage=15", url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("hackernews: build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hackernews: request failed: %w", err)
	}
	defer resp.Body.Close()

	var algoliaResp AlgoliaResponse
	if err := json.NewDecoder(resp.Body).Decode(&algoliaResp); err != nil {
		return nil, fmt.Errorf("hackernews: decode: %w", err)
	}

	var results []CategoryResult
	for _, hit := range algoliaResp.Hits {
		content := hit.StoryText
		if len(content) > 300 {
			content = content[:300] + "..."
		}

		u := hit.URL
		if u == "" {
			u = fmt.Sprintf("https://news.ycombinator.com/item?id=%s", hit.ObjectID)
		}

		results = append(results, CategoryResult{
			Title:       hit.Title,
			URL:         u,
			Content:     content,
			Category:    category,
			Source:      "hackernews",
			PublishedAt: hit.CreatedAt,
		})
	}

	d.Logger.Info(ctx, "hackernews: fetch done",
		logger.Field{Key: "category", Value: category},
		logger.Field{Key: "count", Value: len(results)},
	)

	return results, nil
}
