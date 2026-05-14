package rsshub

type RSSHubConfig struct {
	BaseURL    string `yaml:"base_url"`
	Timeout    int    `yaml:"timeout"`
	MaxRetries int    `yaml:"max_retries"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	RSSHub RSSHubConfig `yaml:"rsshub"`
	Logger LoggerConfig `yaml:"logger"`
}

type Feed struct {
	Title         string `json:"title"`
	Link          string `json:"link"`
	Description   string `json:"description"`
	LastBuildDate string `json:"lastBuildDate"`
	Items         []Item `json:"items"`
}

type Item struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	GUID        string `json:"guid"`
	Description string `json:"description"`
	Author      string `json:"author,omitempty"`
	PubDate     string `json:"pubDate,omitempty"`
	Category    string `json:"category,omitempty"`
}


