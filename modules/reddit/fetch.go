package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chandra179/gosdk/logger"
)

var CategorySubreddits = map[string][]string{
	"economy": {"Economics", "economy"},
	"tech":    {"technology", "programming"},
}

type CategoryResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
}

type hotListing struct {
	Data struct {
		Children []struct {
			Data struct {
				Title     string  `json:"title"`
				URL       string  `json:"url"`
				Selftext  string  `json:"selftext"`
				Permalink string  `json:"permalink"`
				Created   float64 `json:"created_utc"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (d *Dependencies) FetchCategory(category string) ([]CategoryResult, error) {
	ctx := context.Background()

	subs, ok := CategorySubreddits[category]
	if !ok {
		return nil, fmt.Errorf("unknown category: %s (allowed: economy, tech)", category)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	var results []CategoryResult

	for _, sub := range subs {
		url := fmt.Sprintf("%s/r/%s/hot.json?limit=10", d.Config.Reddit.BaseURL, sub)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			d.Logger.Warn(ctx, "reddit: build request failed",
				logger.Field{Key: "sub", Value: sub},
				logger.Field{Key: "err", Value: err},
			)
			continue
		}
		req.Header.Set("User-Agent", d.Config.Reddit.UserAgent)

		resp, err := client.Do(req)
		if err != nil {
			d.Logger.Warn(ctx, "reddit: request failed",
				logger.Field{Key: "sub", Value: sub},
				logger.Field{Key: "err", Value: err},
			)
			continue
		}

		var listing hotListing
		if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
			resp.Body.Close()
			d.Logger.Warn(ctx, "reddit: decode failed",
				logger.Field{Key: "sub", Value: sub},
				logger.Field{Key: "err", Value: err},
			)
			continue
		}
		resp.Body.Close()

		for _, child := range listing.Data.Children {
			u := child.Data.URL
			if u == "" || child.Data.Selftext != "" {
				u = "https://www.reddit.com" + child.Data.Permalink
			}
			content := child.Data.Selftext
			if len(content) > 300 {
				content = content[:300] + "..."
			}
			results = append(results, CategoryResult{
				Title:       child.Data.Title,
				URL:         u,
				Content:     content,
				Category:    category,
				Source:      "reddit",
				PublishedAt: time.Unix(int64(child.Data.Created), 0).Format(time.RFC3339),
			})
		}
	}

	return results, nil
}
