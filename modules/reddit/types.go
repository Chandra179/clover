package reddit

import "time"

type RedditConfig struct {
	BaseURL    string        `yaml:"base_url"`
	UserAgent  string        `yaml:"user_agent"`
	RateLimit  time.Duration `yaml:"rate_limit"`
	MaxRetries int           `yaml:"max_retries"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	Reddit RedditConfig `yaml:"reddit"`
	Logger LoggerConfig `yaml:"logger"`
}

type Listing struct {
	Kind string     `json:"kind"`
	Data ListingData `json:"data"`
}

type ListingData struct {
	After    string    `json:"after"`
	Before   *string   `json:"before"`
	Dist     int       `json:"dist"`
	Children []Thing   `json:"children"`
}

type Thing struct {
	Kind string `json:"kind"`
	Data Post   `json:"data"`
}

type Post struct {
	ID                     string  `json:"id"`
	Name                   string  `json:"name"`
	Title                  string  `json:"title"`
	Selftext               string  `json:"selftext"`
	SelftextHTML           *string `json:"selftext_html"`
	URL                    string  `json:"url"`
	Permalink              string  `json:"permalink"`
	Author                 string  `json:"author"`
	Subreddit              string  `json:"subreddit"`
	SubredditPrefixed      string  `json:"subreddit_name_prefixed"`
	SubredditSubscribers   int     `json:"subreddit_subscribers"`
	Score                  int     `json:"score"`
	Ups                    int     `json:"ups"`
	UpvoteRatio            float64 `json:"upvote_ratio"`
	NumComments            int     `json:"num_comments"`
	CreatedUTC             float64 `json:"created_utc"`
	Over18                 bool    `json:"over_18"`
	Spoiler                bool    `json:"spoiler"`
	Stickied               bool    `json:"stickied"`
	LinkFlairText          *string `json:"link_flair_text"`
	Thumbnail              string  `json:"thumbnail"`
	Domain                 string  `json:"domain"`
	IsSelf                 bool    `json:"is_self"`
	IsVideo                bool    `json:"is_video"`
	NumCrossposts          int     `json:"num_crossposts"`
	TotalAwardsReceived    int     `json:"total_awards_received"`
	AllowLiveComments      bool    `json:"allow_live_comments"`
	IsOriginalContent      bool    `json:"is_original_content"`
	IsRedditMediaDomain    bool    `json:"is_reddit_media_domain"`
	IsCreatedFromAdsUI     bool    `json:"is_created_from_ads_ui"`
	IsCrosspostable        bool    `json:"is_crosspostable"`
	IsMeta                 bool    `json:"is_meta"`
	IsRobotIndexable       bool    `json:"is_robot_indexable"`
	AuthorPremium          bool    `json:"author_premium"`
	Gilded                 int     `json:"gilded"`
	Distinguished          *string `json:"distinguished"`
	Pinned                 bool    `json:"pinned"`
	Locked                 bool    `json:"locked"`
	Archived               bool    `json:"archived"`
	Hidden                 bool    `json:"hidden"`
	Quarantine             bool    `json:"quarantine"`
	Saved                  bool    `json:"saved"`
	Visited                bool    `json:"visited"`
	Clicked                bool    `json:"clicked"`
	NoFollow               bool    `json:"no_follow"`
	SendReplies            bool    `json:"send_replies"`
	ContestMode            bool    `json:"contest_mode"`
}

type CommentListing []interface{}

type Comment struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Author         string  `json:"author"`
	Body           string  `json:"body"`
	BodyHTML       string  `json:"body_html"`
	Score          int     `json:"score"`
	CreatedUTC     float64 `json:"created_utc"`
	ParentID       string  `json:"parent_id"`
	LinkID         string  `json:"link_id"`
	Depth          int     `json:"depth"`
	Stickied       bool    `json:"stickied"`
	Distinguished  *string `json:"distinguished"`
	Subreddit      string  `json:"subreddit"`
	SubredditPrefixed string `json:"subreddit_name_prefixed"`
	TotalAwards    int     `json:"total_awards_received"`
}

type SubredditInfo struct {
	Kind string        `json:"kind"`
	Data SubredditData `json:"data"`
}

type SubredditData struct {
	ID                        string  `json:"id"`
	DisplayName               string  `json:"display_name"`
	DisplayNamePrefixed       string  `json:"display_name_prefixed"`
	Title                     string  `json:"title"`
	PublicDescription         string  `json:"public_description"`
	Subscribers               int     `json:"subscribers"`
	CreatedUTC                float64 `json:"created_utc"`
	Over18                    bool    `json:"over18"`
	Lang                      string  `json:"lang"`
	SubredditType             string  `json:"subreddit_type"`
	UserIsBanned              *bool   `json:"user_is_banned"`
	UserFlairEnabledInSr      *bool   `json:"user_flair_enabled_in_sr"`
}

type SearchParams struct {
	Query      string
	Sort       string
	Timeframe  string
	Limit      int
	RestrictSR bool
	After      string
}
