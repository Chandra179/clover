package handler

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"brook/modules/hackernews"
	"brook/modules/lobsters"
	"brook/modules/reddit"
	"brook/modules/rss"
	"brook/modules/rsshub"
	"brook/modules/wikipedia"
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
	Reddit     *reddit.Dependencies
	Wikipedia  *wikipedia.Dependencies
	HackerNews *hackernews.Dependencies
	Lobsters   *lobsters.Dependencies
	RSS        *rss.Dependencies
	RSSHub     *rsshub.Dependencies
}

func NewDependencies(rd *reddit.Dependencies, wk *wikipedia.Dependencies, hn *hackernews.Dependencies, lb *lobsters.Dependencies, rs *rss.Dependencies, rh *rsshub.Dependencies) *Dependencies {
	return &Dependencies{
		Reddit:     rd,
		Wikipedia:  wk,
		HackerNews: hn,
		Lobsters:   lb,
		RSS:        rs,
		RSSHub:     rh,
	}
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		results, err := d.Reddit.FetchCategory(category)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, "reddit: "+err.Error())
			return
		}
		items := make([]NewsItem, len(results))
		for i, r := range results {
			items[i] = NewsItem{Title: r.Title, URL: r.URL, Content: r.Content, Category: r.Category, Source: r.Source, PublishedAt: r.PublishedAt}
		}
		allItems = append(allItems, items...)
		sources["reddit"] = len(items)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		results, err := d.Wikipedia.FetchCategory(category)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, "wikipedia: "+err.Error())
			return
		}
		items := make([]NewsItem, len(results))
		for i, r := range results {
			items[i] = NewsItem{Title: r.Title, URL: r.URL, Content: r.Content, Category: r.Category, Source: r.Source, PublishedAt: r.PublishedAt}
		}
		allItems = append(allItems, items...)
		sources["wikipedia"] = len(items)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		results, err := d.HackerNews.FetchCategory(category)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, "hackernews: "+err.Error())
			return
		}
		items := make([]NewsItem, len(results))
		for i, r := range results {
			items[i] = NewsItem{Title: r.Title, URL: r.URL, Content: r.Content, Category: r.Category, Source: r.Source, PublishedAt: r.PublishedAt}
		}
		allItems = append(allItems, items...)
		sources["hackernews"] = len(items)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		results, err := d.Lobsters.FetchCategory(category)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, "lobsters: "+err.Error())
			return
		}
		items := make([]NewsItem, len(results))
		for i, r := range results {
			items[i] = NewsItem{Title: r.Title, URL: r.URL, Content: r.Content, Category: r.Category, Source: r.Source, PublishedAt: r.PublishedAt}
		}
		allItems = append(allItems, items...)
		sources["lobsters"] = len(items)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		results, err := d.RSS.FetchCategory(category)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, "rss: "+err.Error())
			return
		}
		items := make([]NewsItem, len(results))
		for i, r := range results {
			items[i] = NewsItem{Title: r.Title, URL: r.URL, Content: r.Content, Category: r.Category, Source: r.Source, PublishedAt: r.PublishedAt}
		}
		allItems = append(allItems, items...)
		sources["rss"] = len(items)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		results, err := d.RSSHub.FetchCategory(category)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, "rsshub: "+err.Error())
			return
		}
		items := make([]NewsItem, len(results))
		for i, r := range results {
			items[i] = NewsItem{Title: r.Title, URL: r.URL, Content: r.Content, Category: r.Category, Source: r.Source, PublishedAt: r.PublishedAt}
		}
		allItems = append(allItems, items...)
		sources["rsshub"] = len(items)
	}()

	wg.Wait()

	seen := map[string]bool{}
	deduped := make([]NewsItem, 0, len(allItems))
	for _, item := range allItems {
		key := strings.ToLower(strings.Join(strings.Fields(item.Title), " "))
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, item)
	}

	resp := NewsResponse{
		Category: category,
		Count:    len(deduped),
		Sources:  sources,
		Items:    deduped,
	}

	if len(errs) > 0 && len(deduped) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "all sources failed",
			"detail": errs,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
