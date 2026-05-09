package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Chandra179/gosdk/logger"
	"github.com/gin-gonic/gin"

	"brook/handler"
	"brook/modules/reddit"
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
		Logger: reddit.LoggerConfig{
			Level: "dev",
		},
	})

	wk := wikipedia.NewDependencies(wikipedia.Config{
		Wikipedia: wikipedia.WikiConfig{
			BaseURL:   "https://en.wikipedia.org",
			UserAgent: "brook/1.0",
		},
		Logger: wikipedia.LoggerConfig{
			Level: "dev",
		},
	})

	h := handler.NewDependencies(rd, wk)

	r := gin.Default()
	r.GET("/news", h.NewsHandler)

	log.Info(ctx, "starting server", logger.Field{Key: "addr", Value: ":8080"})
	if err := r.Run(":8080"); err != nil {
		log.Error(ctx, "server failed", logger.Field{Key: "err", Value: err})
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
