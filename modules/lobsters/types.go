package lobsters

type LobstersConfig struct {
	BaseURL    string `yaml:"base_url"`
	MaxRetries int    `yaml:"max_retries"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	Lobsters LobstersConfig `yaml:"lobsters"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type Story struct {
	ShortID      string      `json:"short_id"`
	ShortIDURL   string      `json:"short_id_url"`
	CreatedAt    string      `json:"created_at"`
	Title        string      `json:"title"`
	URL          string      `json:"url"`
	Score        int         `json:"score"`
	CommentCount int         `json:"comment_count"`
	Description  string      `json:"description"`
	CommentsURL  string      `json:"comments_url"`
	SubmitterUser UserRef    `json:"submitter_user"`
	Tags         []string    `json:"tags"`
	Comments     []Comment   `json:"comments,omitempty"`
}

type UserRef struct {
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	IsAdmin   bool   `json:"is_admin"`
}

type Comment struct {
	ShortID    string `json:"short_id"`
	ShortIDURL string `json:"short_id_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	IsDeleted  bool   `json:"is_deleted"`
	IsModerated bool  `json:"is_moderated"`
	Score      int    `json:"score"`
	Comment    string `json:"comment"`
	IndentLevel int   `json:"indent_level"`
	Hat        *string `json:"hat"`
	IssuerID   *string `json:"issuer_id"`
}

type Tag struct {
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
	HotnessMod  float64 `json:"hotness_mod"`
	IsMedia     bool    `json:"is_media"`
}

type UserDetail struct {
	Username  string   `json:"username"`
	CreatedAt string   `json:"created_at"`
	IsAdmin   bool     `json:"is_admin"`
	About     string   `json:"about"`
	InvitedBy string   `json:"invited_by,omitempty"`
	Stories   []Story  `json:"stories"`
	Comments  []Comment `json:"comments"`
}

type CategoryResult struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Source   string `json:"source"`
}
