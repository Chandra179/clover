package rsshub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chandra179/gosdk/logger"
)

var CategoryRoutes = map[string]string{
	"economy": "reuters/topics/business",
	"tech":    "github/trending/daily",
}

func (d *Dependencies) FetchCategory(category string) ([]CategoryResult, error) {
	ctx := context.Background()

	route, ok := CategoryRoutes[category]
	if !ok {
		return nil, fmt.Errorf("unknown category: %s (allowed: economy, tech)", category)
	}

	client := &http.Client{Timeout: 15 * time.Second}

	url := fmt.Sprintf("%s/%s.json", d.Config.RSSHub.BaseURL, route)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("rsshub: build request: %w", err)
	}
	req.Header.Set("User-Agent", "brook/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rsshub: request failed: %w", err)
	}
	defer resp.Body.Close()

	var feed Feed
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("rsshub: decode: %w", err)
	}

	limit := 10
	if len(feed.Items) > limit {
		feed.Items = feed.Items[:limit]
	}

	var results []CategoryResult
	for _, item := range feed.Items {
		content := item.Description
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		results = append(results, CategoryResult{
			Title:       item.Title,
			URL:         item.Link,
			Content:     content,
			Category:    category,
			Source:      "rsshub",
			PublishedAt: item.PubDate,
		})
	}

	d.Logger.Info(ctx, "rsshub: fetch done",
		logger.Field{Key: "category", Value: category},
		logger.Field{Key: "route", Value: route},
		logger.Field{Key: "count", Value: len(results)},
	)

	return results, nil
}
