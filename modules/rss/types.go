package rss

import "encoding/xml"

type RSSConfig struct {
	BaseURL    string              `yaml:"base_url"`
	Categories map[string][]string `yaml:"categories"`
	MaxRetries int                 `yaml:"max_retries"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	RSS    RSSConfig    `yaml:"rss"`
	Logger LoggerConfig `yaml:"logger"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Author      string `xml:"author"`
	Category    string `xml:"category"`
}

type AtomFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Title   string     `xml:"title"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	Title   string    `xml:"title"`
	Link    AtomLink  `xml:"link"`
	ID      string    `xml:"id"`
	Published string  `xml:"published"`
	Updated string    `xml:"updated"`
	Summary string    `xml:"summary"`
	Content string    `xml:"content"`
	Author  AtomAuthor `xml:"author"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
}


