package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Chandra179/gosdk/logger"
)

func (d *Dependencies) FetchCategory(category string) ([]CategoryResult, error) {
	ctx := context.Background()

	feedURLs, ok := d.Config.RSS.Categories[category]
	if !ok || len(feedURLs) == 0 {
		return nil, fmt.Errorf("rss: no feeds configured for category: %s", category)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	var results []CategoryResult

	for _, feedURL := range feedURLs {
		items, err := d.fetchFeed(ctx, client, feedURL)
		if err != nil {
			d.Logger.Warn(ctx, "rss: feed fetch failed",
				logger.Field{Key: "url", Value: feedURL},
				logger.Field{Key: "err", Value: err},
			)
			continue
		}
		for _, item := range items {
			results = append(results, CategoryResult{
				Title:    item.Title,
				URL:      item.Link,
				Content:  truncate(item.Content, 300),
				Category: category,
				Source:   "rss",
			})
		}
	}

	d.Logger.Info(ctx, "rss: fetch done",
		logger.Field{Key: "category", Value: category},
		logger.Field{Key: "feed_count", Value: len(feedURLs)},
		logger.Field{Key: "total", Value: len(results)},
	)

	return results, nil
}

func (d *Dependencies) fetchFeed(ctx context.Context, client *http.Client, feedURL string) ([]rawItem, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("rss: build request: %w", err)
	}
	req.Header.Set("User-Agent", "brook/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rss: fetch %s: %w", feedURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rss: read body: %w", err)
	}

	return d.parseFeed(body)
}

type rawItem struct {
	Title   string
	Link    string
	Content string
}

func (d *Dependencies) parseFeed(body []byte) ([]rawItem, error) {
	trimmed := strings.TrimSpace(string(body))

	if strings.HasPrefix(trimmed, "<?xml") {
		if idx := strings.Index(trimmed, "<rss"); idx >= 0 {
			trimmed = trimmed[idx:]
		} else if idx := strings.Index(trimmed, "<feed"); idx >= 0 {
			trimmed = trimmed[idx:]
		} else if idx := strings.Index(trimmed, "<rdf:"); idx >= 0 {
			trimmed = trimmed[idx:]
		}
	}

	var rssFeed RSS
	if err := xml.Unmarshal([]byte(trimmed), &rssFeed); err == nil && len(rssFeed.Channel.Items) > 0 {
		var items []rawItem
		for _, item := range rssFeed.Channel.Items {
			content := item.Description
			if content == "" {
				content = item.Title
			}
			link := item.Link
			if link == "" {
				link = item.GUID
			}
			items = append(items, rawItem{
				Title:   item.Title,
				Link:    link,
				Content: content,
			})
		}
		return items, nil
	}

	var atomFeed AtomFeed
	if err := xml.Unmarshal([]byte(trimmed), &atomFeed); err == nil && len(atomFeed.Entries) > 0 {
		var items []rawItem
		for _, entry := range atomFeed.Entries {
			content := entry.Summary
			if content == "" {
				content = entry.Content
			}
			if content == "" {
				content = entry.Title
			}
			items = append(items, rawItem{
				Title:   entry.Title,
				Link:    entry.Link.Href,
				Content: content,
			})
		}
		return items, nil
	}

	return nil, fmt.Errorf("rss: unrecognized feed format")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
