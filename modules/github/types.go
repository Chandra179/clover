package github

type GHConfig struct {
	BaseURL    string `yaml:"base_url"`
	MaxRetries int    `yaml:"max_retries"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	GitHub GHConfig     `yaml:"github"`
	Logger LoggerConfig `yaml:"logger"`
}

type SearchResponse struct {
	Items      []Repo `json:"items"`
	TotalCount int    `json:"total_count"`
}

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	Topics      []string `json:"topics"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Owner       Owner  `json:"owner"`
}

type Owner struct {
	Login string `json:"login"`
}
