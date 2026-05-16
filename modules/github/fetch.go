package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"brook/modules/common"

	"github.com/Chandra179/gosdk/logger"
)

var CategorySearchQueries = map[string]string{
	"economy":  "fintech OR blockchain OR finance",
	"tech":     "programming OR technology OR software",
	"science":  "science OR research OR scientific",
	"ai":       "artificial-intelligence OR machine-learning OR deep-learning",
	"security": "security OR cybersecurity OR privacy",
	"startups": "startup OR startups OR YCombinator",
}

func (d *Dependencies) FetchCategory(category string) ([]common.CategoryResult, error) {
	ctx := context.Background()

	query, ok := CategorySearchQueries[category]
	if !ok {
		return nil, fmt.Errorf("github: unknown category: %s", category)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	apiURL := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=15",
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("github: build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "brook/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			body = map[string]any{"raw": "decode failed"}
		}
		d.Logger.Error(ctx, "github: API error",
			logger.Field{Key: "status", Value: resp.StatusCode},
			logger.Field{Key: "body", Value: body},
		)
		return nil, fmt.Errorf("github: API returned %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("github: decode: %w", err)
	}

	var results []common.CategoryResult
	for _, repo := range searchResp.Items {
		content := repo.Description
		if content == "" {
			content = fmt.Sprintf("A %s repository by %s", repo.Language, repo.Owner.Login)
		}
		if len(content) > 300 {
			content = content[:300] + "..."
		}

		results = append(results, common.CategoryResult{
			Title:       repo.Name,
			URL:         repo.HTMLURL,
			Content:     content,
			Category:    category,
			Source:      "github",
			PublishedAt: repo.UpdatedAt,
		})
	}

	d.Logger.Info(ctx, "github: fetch done",
		logger.Field{Key: "category", Value: category},
		logger.Field{Key: "count", Value: len(results)},
	)

	return results, nil
}
