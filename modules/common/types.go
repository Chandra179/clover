package common

type CategoryResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
}

type Fetcher interface {
	FetchCategory(category string) ([]CategoryResult, error)
}
