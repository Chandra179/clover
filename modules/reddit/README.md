# Reddit Data Collector

Free JSON API for scraping Reddit public data. No API key required for read-only access.

## Base URL

```
https://www.reddit.com
```

## Authentication

| Tier | Method | Rate Limit |
|------|--------|------------|
| No auth | Append `.json` to any URL | 60 req/min per IP |
| OAuth2 (free) | `Authorization: Bearer <token>` | 600 req/10min |

No-auth is sufficient for all read-only listings below.

## Endpoints

### Subreddit Listings (no auth)

Listings are paginated collections of posts. All share common pagination params.

| Endpoint | Description |
|----------|-------------|
| `GET /r/{subreddit}/hot.json` | Hot posts (default sort) |
| `GET /r/{subreddit}/new.json` | Newest posts |
| `GET /r/{subreddit}/top.json` | Top posts (by `t` timeframe) |
| `GET /r/{subreddit}/rising.json` | Rising posts |
| `GET /r/{subreddit}/controversial.json` | Controversial posts |
| `GET /r/{subreddit}/random.json` | Random post (redirects) |

### Front Page / All (no auth)

| Endpoint | Description |
|----------|-------------|
| `GET /hot.json` | Front page hot |
| `GET /new.json` | Front page new |
| `GET /r/all/hot.json` | r/all hot |

### Subreddit Info (no auth)

| Endpoint | Description |
|----------|-------------|
| `GET /r/{subreddit}/about.json` | Subreddit metadata (subscribers, rules, etc.) |
| `GET /r/{subreddit}/about/rules.json` | Subreddit rules |
| `GET /r/{subreddit}/about/moderators.json` | Moderator list |

### Comments (no auth)

| Endpoint | Description |
|----------|-------------|
| `GET /comments/{post_id}.json` | Comments on a post (post_id = article id, e.g. `1t3r6ax`) |
| `GET /r/{subreddit}/comments/{post_id}/{slug}.json` | Comments with subreddit context |

### Posts by ID (no auth)

| Endpoint | Description |
|----------|-------------|
| `GET /api/info.json?id=t3_{id}` | Fetch post(s) by fullname |
| `GET /by_id/t3_{id1},t3_{id2}.json` | Fetch multiple posts by fullname |

### Search (no auth)

| Endpoint | Description |
|----------|-------------|
| `GET /r/{subreddit}/search.json?q={query}` | Search within subreddit |
| `GET /search.json?q={query}` | Site-wide search |

### User Data (no auth)

| Endpoint | Description |
|----------|-------------|
| `GET /user/{username}/about.json` | User profile info |
| `GET /user/{username}/submitted.json` | User's posts |
| `GET /user/{username}/comments.json` | User's comments |
| `GET /user/{username}/overview.json` | User's overview (posts + comments) |

### Subreddit Discovery (no auth)

| Endpoint | Description |
|----------|-------------|
| `GET /subreddits/popular.json` | Popular subreddits |
| `GET /subreddits/new.json` | New subreddits |
| `GET /subreddits/search.json?q={query}` | Search subreddits |
| `GET /api/subreddit_autocomplete.json?query={q}` | Autocomplete subreddit names |

## Common Pagination Parameters (Listings)

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `limit` | int | 25 | Items per page (1–100) |
| `after` | string | — | Fullname of item after which to fetch (e.g. `t3_1t3r6ax`) |
| `before` | string | — | Fullname of item before which to fetch |
| `count` | int | 0 | Items already seen in this listing |
| `show` | string | — | Pass `all` to disable "hide voted" filter |

## Listing-Specific Parameters

### `/top.json` & `/controversial.json`

| Param | Values | Default | Description |
|-------|--------|---------|-------------|
| `t` | `hour`, `day`, `week`, `month`, `year`, `all` | `day` | Timeframe |

### Comments (`/comments/{post_id}.json`)

| Param | Values | Default | Description |
|-------|--------|---------|-------------|
| `sort` | `confidence`, `top`, `new`, `controversial`, `old`, `qa` | `confidence` | Comment sort order |
| `limit` | int | — | Max comments to return |
| `depth` | int | — | Max comment tree depth |
| `threaded` | bool | `true` | Return threaded vs flat |

### Search (`/search.json`)

| Param | Type | Description |
|-------|------|-------------|
| `q` | string | Search query |
| `sort` | `relevance`, `hot`, `new`, `top`, `comments` | Sort order |
| `t` | `hour`, `day`, `week`, `month`, `year`, `all` | Timeframe |
| `limit` | int | Results per page |
| `restrict_sr` | bool | Restrict to current subreddit |
| `type` | `link`, `sr`, `user` | Filter by type |

## Rate Limits

| Auth | Limit | Window |
|------|-------|--------|
| No auth | 60 requests | Per minute, per IP |
| OAuth2 | 600 requests | Per 10 minutes |
| OAuth2 (registered app) | 6000 requests | Per 10 minutes |

**Important**: Reddit returns HTTP 429 with `Retry-After` header when rate-limited. Respect it.

Headers:
- `X-Ratelimit-Used`: Approx number used in current period
- `X-Ratelimit-Remaining`: Approx number remaining
- `X-Ratelimit-Reset`: Approx seconds until reset

## Response Format

Reddit wraps all responses in a type-tagged JSON structure:

```json
{
  "kind": "Listing",
  "data": {
    "after": "t3_abc123",
    "before": null,
    "dist": 25,
    "children": [
      {
        "kind": "t3",
        "data": { /* ...post fields... */ }
      }
    ]
  }
}
```

### Thing Types (`kind`)

| `kind` | Type | Description |
|--------|------|-------------|
| `t1_` | Comment | A comment |
| `t2_` | Account | A user account |
| `t3_` | Link | A post/link |
| `t4_` | Message | A private message |
| `t5_` | Subreddit | A subreddit |
| `t6_` | Award | An award |
| `Listing` | — | Wrapper for paginated lists |

### Key Post Fields (in `data`)

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Post ID (base-36) |
| `name` | string | Fullname (`t3_` + id) |
| `title` | string | Post title |
| `selftext` | string | Body text (self posts) or empty |
| `selftext_html` | string | Rendered HTML body |
| `url` | string | Link URL or Reddit permalink |
| `permalink` | string | Relative permalink |
| `author` | string | Author username |
| `subreddit` | string | Subreddit name |
| `subreddit_name_prefixed` | string | e.g. `r/golang` |
| `subreddit_subscribers` | int | Subreddit subscriber count |
| `score` | int | Net score (upvotes - downvotes) |
| `ups` | int | Upvote count (fuzzed for old posts) |
| `upvote_ratio` | float | Percentage upvoted |
| `num_comments` | int | Number of comments |
| `created_utc` | float | Unix timestamp (UTC) |
| `created` | float | Unix timestamp (local) |
| `over_18` | bool | NSFW flag |
| `spoiler` | bool | Spoiler flag |
| `stickied` | bool | Stickied by mod |
| `link_flair_text` | string | Flair text |
| `thumbnail` | string | Thumbnail URL or `"self"`, `"default"`, `"nsfw"` |
| `domain` | string | Link domain |
| `is_self` | bool | Self-post (text only) |
| `is_video` | bool | Video post |
| `media` | object | Media embed data (null for no media) |
| `preview` | object | Image preview data |
| `gallery_data` | object | Gallery data (multi-image posts) |
| `poll_data` | object | Poll data |

### Key Comment Fields (in `data`)

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Comment ID |
| `name` | string | Fullname (`t1_` + id) |
| `author` | string | Author username |
| `body` | string | Comment body (markdown) |
| `body_html` | string | Rendered HTML |
| `score` | int | Net score |
| `created_utc` | float | Unix timestamp |
| `parent_id` | string | Fullname of parent (`t3_` for post, `t1_` for comment) |
| `link_id` | string | Fullname of the post this comment belongs to |
| `replies` | object/string | Nested comment tree or `""` if none |
| `depth` | int | Depth in tree |
| `stickied` | bool | Stickied by mod |
| `distinguished` | string | `"moderator"` or `"admin"` |

### Key Subreddit Fields (about.json)

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Subreddit ID |
| `display_name` | string | Subreddit name |
| `title` | string | Display title |
| `public_description` | string | Sidebar description |
| `subscribers` | int | Subscriber count |
| `created_utc` | float | Creation timestamp |
| `over18` | bool | NSFW flag |
| `lang` | string | Language code |

### Raw JSON Encoding Note

By default, response bodies HTML-encode `<`, `>`, and `&`. Pass `raw_json=1` param to get raw JSON.

## Example Usage

```bash
# Fetch hot posts from r/golang
curl -H "User-Agent: clover/0.1" \
  "https://www.reddit.com/r/golang/hot.json?limit=10&raw_json=1"

# Search for "monad" in r/golang
curl -H "User-Agent: clover/0.1" \
  "https://www.reddit.com/r/golang/search.json?q=monad&limit=10&restrict_sr=1&raw_json=1"

# Get comments on a post
curl -H "User-Agent: clover/0.1" \
  "https://www.reddit.com/r/golang/comments/1t3r6ax.json?raw_json=1"

# Get subreddit info
curl -H "User-Agent: clover/0.1" \
  "https://www.reddit.com/r/golang/about.json?raw_json=1"

# Pagination — use `after` from previous response
curl -H "User-Agent: clover/0.1" \
  "https://www.reddit.com/r/golang/hot.json?limit=10&after=t3_1t3r6ax&raw_json=1"
```

**User-Agent header is required.** Reddit blocks requests without it. Format: `<app>/<version>`.

## Limitations

- Max 100 items per listing request (`limit` parameter)
- Max 1000 posts accessible per listing (Reddit's internal listing cache limit)
- Deleted/removed posts return minimal data (`selftext: "[removed]"` or `"[deleted]"`)
- Score values are fuzzed to prevent spam
- No historical search beyond listing cache (use Pushshift.io for archival access)
