package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"brook/modules/reddit"
	"brook/modules/wikipedia"
)

type NewsItem struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Source   string `json:"source"`
}

type NewsResponse struct {
	Category string     `json:"category"`
	Count    int        `json:"count"`
	Items    []NewsItem `json:"items"`
}

type Dependencies struct {
	Reddit    *reddit.Dependencies
	Wikipedia *wikipedia.Dependencies
}

func NewDependencies(rd *reddit.Dependencies, wk *wikipedia.Dependencies) *Dependencies {
	return &Dependencies{Reddit: rd, Wikipedia: wk}
}

func (d *Dependencies) NewsHandler(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		category = "tech"
	}

	var (
		redditResults []reddit.CategoryResult
		wg            sync.WaitGroup
		mu            sync.Mutex
		errs          []string
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		results, err := d.Reddit.FetchCategory(category)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, "reddit: "+err.Error())
			return
		}
		redditResults = results
	}()

	wg.Wait()

	items := make([]NewsItem, 0, len(redditResults))
	for _, r := range redditResults {
		items = append(items, NewsItem{
			Title:    r.Title,
			URL:      r.URL,
			Content:  r.Content,
			Category: r.Category,
			Source:   r.Source,
		})
	}

	resp := NewsResponse{
		Category: category,
		Count:    len(items),
		Items:    items,
	}

	if len(errs) > 0 && len(items) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "all sources failed",
			"detail": errs,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
