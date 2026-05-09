package wikipedia

type WikiConfig struct {
	BaseURL    string `yaml:"base_url"`
	UserAgent  string `yaml:"user_agent"`
	MaxLag     int    `yaml:"max_lag"`
	MaxRetries int    `yaml:"max_retries"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	Wikipedia WikiConfig   `yaml:"wikipedia"`
	Logger    LoggerConfig `yaml:"logger"`
}

type QueryResponse struct {
	Batchcomplete string          `json:"batchcomplete"`
	Continue      *ContinueBlock  `json:"continue"`
	Query         *QueryResult    `json:"query"`
	Error         *APIError       `json:"error"`
}

type ContinueBlock map[string]string

type QueryResult struct {
	Normalized []NormalizedTitle   `json:"normalized"`
	Pages      map[string]Page     `json:"pages"`
	Search     []SearchResult      `json:"search"`
	Random     []RandomResult      `json:"random"`
	SearchInfo *SearchInfo         `json:"searchinfo"`
}

type Page struct {
	PageID      int                `json:"pageid"`
	NS          int                `json:"ns"`
	Title       string             `json:"title"`
	Extract     string             `json:"extract"`
	Revisions   []Revision         `json:"revisions"`
	Categories  []CategoryEntry    `json:"categories"`
	Links       []LinkEntry        `json:"links"`
	LangLinks   []LangLinkEntry    `json:"langlinks"`
	Templates   []TemplateEntry    `json:"templates"`
	ExtLinks    []ExtLinkEntry     `json:"extlinks"`
	Images      []ImageEntry       `json:"images"`
	Thumbnail   *Thumbnail         `json:"thumbnail"`
	PageImage   string             `json:"pageimage"`
	Timestamp   string             `json:"timestamp"`
	Touched     string             `json:"touched"`
	Length      int                `json:"length"`
	Redirect    bool               `json:"redirect"`
	FullURL     string             `json:"fullurl"`
	CanonicalURL string            `json:"canonicalurl"`
	ContentModel string            `json:"contentmodel"`
	PageLanguage string            `json:"pagelanguage"`
	PageProps   map[string]string  `json:"pageprops"`
	Coordinates []Coordinate       `json:"coordinates"`
}

type Revision struct {
	RevID         int         `json:"revid"`
	ParentID      int         `json:"parentid"`
	Timestamp     string      `json:"timestamp"`
	User          string      `json:"user"`
	Comment       string      `json:"comment"`
	Slots         *struct {
		Main *struct {
			ContentModel  string `json:"contentmodel"`
			ContentFormat string `json:"contentformat"`
			Content       string `json:"content"`
		} `json:"main"`
	} `json:"slots"`
}

type SearchResult struct {
	NS        int    `json:"ns"`
	Title     string `json:"title"`
	PageID    int    `json:"pageid"`
	Size      int    `json:"size"`
	WordCount int    `json:"wordcount"`
	Snippet   string `json:"snippet"`
	Timestamp string `json:"timestamp"`
}

type SearchInfo struct {
	TotalHits      int    `json:"totalhits"`
	Suggestion     string `json:"suggestion"`
	SuggestionSnippet string `json:"suggestionsnippet"`
}

type RandomResult struct {
	ID    int    `json:"id"`
	NS    int    `json:"ns"`
	Title string `json:"title"`
}

type CategoryEntry struct {
	NS    int    `json:"ns"`
	Title string `json:"title"`
}

type LinkEntry struct {
	NS    int    `json:"ns"`
	Title string `json:"title"`
}

type LangLinkEntry struct {
	Lang        string `json:"lang"`
	Title       string `json:"title"`
	Autonym     string `json:"autonym"`
}

type TemplateEntry struct {
	NS    int    `json:"ns"`
	Title string `json:"title"`
}

type ExtLinkEntry struct {
	URL string `json:"url"`
}

type ImageEntry struct {
	NS    int    `json:"ns"`
	Title string `json:"title"`
}

type Thumbnail struct {
	Source string `json:"source"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Coordinate struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Primary bool    `json:"primary"`
}

type NormalizedTitle struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ParseResponse struct {
	Parse *ParseResult `json:"parse"`
	Error *APIError    `json:"error"`
}

type ParseResult struct {
	Title      string         `json:"title"`
	PageID     int            `json:"pageid"`
	RevID      int            `json:"revid"`
	Text       map[string]interface{} `json:"text"`
	Wikitext   map[string]interface{} `json:"wikitext"`
	Categories []CategoryEntry `json:"categories"`
	Links      []LinkEntry     `json:"links"`
	Templates  []TemplateEntry `json:"templates"`
	Images     []string        `json:"images"`
	Sections   []Section       `json:"sections"`
}

type Section struct {
	TocLevel int    `json:"toclevel"`
	Level    string `json:"level"`
	Line     string `json:"line"`
	Number   string `json:"number"`
	Index    string `json:"index"`
	FromTitle string `json:"fromtitle"`
	ByteOffset int  `json:"byteoffset"`
	Anchor   string `json:"anchor"`
}

type OpenSearchResponse []interface{}

type APIError struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

type SearchParams struct {
	Query     string
	Limit     int
	Namespace int
	Offset    int
}

type PageParams struct {
	Titles       string
	PageIDs      string
	Props        string
	ExIntro      bool
	ExplainText  bool
	ExChars      int
	InProp       string
	ImgThumbSize int
}
