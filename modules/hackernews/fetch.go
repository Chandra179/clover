package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	client := &http.Client{Timeout: 10 * time.Second}

	url := d.Config.HackerNews.BaseURL + "/topstories.json"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("hackernews: build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hackernews: request failed: %w", err)
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("hackernews: decode ids: %w", err)
	}

	limit := 15
	if len(ids) > limit {
		ids = ids[:limit]
	}

	var results []CategoryResult
	for _, id := range ids {
		itemURL := d.Config.HackerNews.BaseURL + fmt.Sprintf("/item/%d.json", id)
		req, err := http.NewRequestWithContext(ctx, "GET", itemURL, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		var item Item
		if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if item.Deleted || item.Dead || item.Title == "" {
			continue
		}

		content := item.Text
		if len(content) > 300 {
			content = content[:300] + "..."
		}

		url := item.URL
		if url == "" {
			url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID)
		}

		results = append(results, CategoryResult{
			Title:    item.Title,
			URL:      url,
			Content:  content,
			Category: category,
			Source:   "hackernews",
		})
	}

	d.Logger.Info(ctx, "hackernews: fetch done",
		logger.Field{Key: "category", Value: category},
		logger.Field{Key: "count", Value: len(results)},
	)

	return results, nil
}
