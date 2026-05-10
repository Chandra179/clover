package hackernews

type HNConfig struct {
	BaseURL    string `yaml:"base_url"`
	MaxRetries int    `yaml:"max_retries"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	HackerNews HNConfig     `yaml:"hackernews"`
	Logger     LoggerConfig `yaml:"logger"`
}

type Item struct {
	ID          int    `json:"id"`
	Deleted     bool   `json:"deleted,omitempty"`
	Type        string `json:"type,omitempty"`
	By          string `json:"by,omitempty"`
	Time        int    `json:"time,omitempty"`
	Text        string `json:"text,omitempty"`
	Dead        bool   `json:"dead,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	Poll        int    `json:"poll,omitempty"`
	Kids        []int  `json:"kids,omitempty"`
	URL         string `json:"url,omitempty"`
	Score       int    `json:"score,omitempty"`
	Title       string `json:"title,omitempty"`
	Parts       []int  `json:"parts,omitempty"`
	Descendants int    `json:"descendants,omitempty"`
}

type User struct {
	ID        string `json:"id"`
	Delay     int    `json:"delay,omitempty"`
	Created   int    `json:"created"`
	Karma     int    `json:"karma"`
	About     string `json:"about,omitempty"`
	Submitted []int  `json:"submitted,omitempty"`
}

type Updates struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}

type CategoryResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
}

type AlgoliaResponse struct {
	Hits        []AlgoliaHit `json:"hits"`
	NbHits      int          `json:"nbHits"`
	Page        int          `json:"page"`
	NbPages     int          `json:"nbPages"`
	HitsPerPage int          `json:"hitsPerPage"`
}

type AlgoliaHit struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Author      string `json:"author"`
	Points      int    `json:"points"`
	StoryText   string `json:"story_text"`
	ObjectID    string `json:"objectID"`
	CreatedAt   string `json:"created_at"`
}
