package handler

import (
	"net/http"
	"sync"

	"brook/modules/common"

	"github.com/gin-gonic/gin"
)

type NewsItem struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
}

type NewsResponse struct {
	Category string         `json:"category"`
	Count    int            `json:"count"`
	Sources  map[string]int `json:"sources"`
	Items    []NewsItem     `json:"items"`
}

type Dependencies struct {
	Fetchers []common.Fetcher
}

func NewDependencies(fetchers []common.Fetcher) *Dependencies {
	return &Dependencies{Fetchers: fetchers}
}

func (d *Dependencies) NewsHandler(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		category = "tech"
	}

	var (
		mu       sync.Mutex
		errs     []string
		wg       sync.WaitGroup
		allItems []NewsItem
		sources  = make(map[string]int)
	)

	for _, f := range d.Fetchers {
		wg.Add(1)
		go func(fetcher common.Fetcher) {
			defer wg.Done()
			results, err := fetcher.FetchCategory(category)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err.Error())
				return
			}
			items := make([]NewsItem, len(results))
			for i, r := range results {
				items[i] = NewsItem{Title: r.Title, URL: r.URL, Content: r.Content, Category: r.Category, Source: r.Source, PublishedAt: r.PublishedAt}
			}
			allItems = append(allItems, items...)
			if len(results) > 0 {
				sources[results[0].Source] = len(items)
			}
		}(f)
	}

	wg.Wait()

	resp := NewsResponse{
		Category: category,
		Count:    len(allItems),
		Sources:  sources,
		Items:    allItems,
	}

	if len(errs) > 0 && len(allItems) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "all sources failed",
			"detail": errs,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
