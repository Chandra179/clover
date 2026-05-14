package lobsters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"brook/modules/common"

	"github.com/Chandra179/gosdk/logger"
)

func (d *Dependencies) FetchCategory(category string) ([]common.CategoryResult, error) {
	ctx := context.Background()

	tagMap := map[string]string{
		"economy": "finance",
		"tech":    "programming",
	}

	tag := tagMap[category]
	if tag == "" && category != "" {
		tag = category
	}

	client := &http.Client{Timeout: 10 * time.Second}

	var url string
	if tag != "" {
		url = fmt.Sprintf("%s/t/%s.json", d.Config.Lobsters.BaseURL, tag)
	} else {
		url = d.Config.Lobsters.BaseURL + "/newest.json"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("lobsters: build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lobsters: request failed: %w", err)
	}
	defer resp.Body.Close()

	var stories []Story
	if err := json.NewDecoder(resp.Body).Decode(&stories); err != nil {
		return nil, fmt.Errorf("lobsters: decode: %w", err)
	}

	limit := 10
	if len(stories) > limit {
		stories = stories[:limit]
	}

	var results []common.CategoryResult
	for _, s := range stories {
		cat := category
		if cat == "" {
			cat = "general"
		}
		results = append(results, common.CategoryResult{
			Title:       s.Title,
			URL:         s.URL,
			Content:     s.Description,
			Category:    cat,
			Source:      "lobsters",
			PublishedAt: s.CreatedAt,
		})
	}

	d.Logger.Info(ctx, "lobsters: fetch done",
		logger.Field{Key: "category", Value: category},
		logger.Field{Key: "count", Value: len(results)},
	)

	return results, nil
}
