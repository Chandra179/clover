package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Chandra179/gosdk/logger"
	"github.com/gin-gonic/gin"

	"brook/handler"
	"brook/modules/common"
	"brook/modules/github"
	"brook/modules/hackernews"
	"brook/modules/lobsters"
	"brook/modules/reddit"
	"brook/modules/rss"
	"brook/modules/rsshub"
	"brook/modules/wikipedia"
)

func main() {
	log := logger.NewLogger("dev")
	ctx := context.Background()

	rd := reddit.NewDependencies(reddit.Config{
		Reddit: reddit.RedditConfig{
			BaseURL:   "https://www.reddit.com",
			UserAgent: "brook/1.0",
		},
		Logger: reddit.LoggerConfig{Level: "dev"},
	})

	wk := wikipedia.NewDependencies(wikipedia.Config{
		Wikipedia: wikipedia.WikiConfig{
			BaseURL:   "https://en.wikipedia.org",
			UserAgent: "brook/1.0",
		},
		Logger: wikipedia.LoggerConfig{Level: "dev"},
	})

	hn := hackernews.NewDependencies(hackernews.Config{
		HackerNews: hackernews.HNConfig{
			BaseURL: "https://hn.algolia.com",
		},
		Logger: hackernews.LoggerConfig{Level: "dev"},
	})

	lb := lobsters.NewDependencies(lobsters.Config{
		Lobsters: lobsters.LobstersConfig{
			BaseURL: "https://lobste.rs",
		},
		Logger: lobsters.LoggerConfig{Level: "dev"},
	})

	rs := rss.NewDependencies(rss.Config{
		RSS: rss.RSSConfig{
			Categories: map[string][]string{
				"economy": {
					"https://feeds.bbci.co.uk/news/business/rss.xml",
					"https://www.cnbc.com/id/100003114/device/rss/rss.html",
				},
				"tech": {
					"https://feeds.bbci.co.uk/news/technology/rss.xml",
					"https://www.theverge.com/rss/index.xml",
				},
				"science": {
					"https://feeds.bbci.co.uk/news/science_and_environment/rss.xml",
					"https://www.nature.com/nature.rss",
				},
				"ai": {
					"https://feeds.feedburner.com/ArtificialIntelligenceNews",
				},
				"security": {
					"https://feeds.bbci.co.uk/news/technology/rss.xml",
				},
				"startups": {
					"https://feeds.feedburner.com/TechCrunch",
				},
			},
		},
		Logger: rss.LoggerConfig{Level: "dev"},
	})

	rh := rsshub.NewDependencies(rsshub.Config{
		RSSHub: rsshub.RSSHubConfig{
			BaseURL: "https://rsshub.app",
		},
		Logger: rsshub.LoggerConfig{Level: "dev"},
	})

	gh := github.NewDependencies(github.Config{
		GitHub: github.GHConfig{
			BaseURL: "https://api.github.com",
		},
		Logger: github.LoggerConfig{Level: "dev"},
	})

	cache := common.NewCache(60 * time.Second)

	fetchers := []common.Fetcher{
		&common.CachedFetcher{Source: "reddit", Fetcher: rd, Cache: cache},
		&common.CachedFetcher{Source: "wikipedia", Fetcher: wk, Cache: cache},
		&common.CachedFetcher{Source: "hackernews", Fetcher: hn, Cache: cache},
		&common.CachedFetcher{Source: "lobsters", Fetcher: lb, Cache: cache},
		&common.CachedFetcher{Source: "rss", Fetcher: rs, Cache: cache},
		&common.CachedFetcher{Source: "rsshub", Fetcher: rh, Cache: cache},
		&common.CachedFetcher{Source: "github", Fetcher: gh, Cache: cache},
	}

	h := handler.NewDependencies(fetchers)

	r := gin.Default()
	r.GET("/news", h.NewsHandler)

	log.Info(ctx, "starting server", logger.Field{Key: "addr", Value: ":8001"})
	if err := r.Run(":8001"); err != nil {
		log.Error(ctx, "server failed", logger.Field{Key: "err", Value: err})
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
