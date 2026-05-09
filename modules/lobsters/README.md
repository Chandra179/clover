# Lobste.rs Data Collector

Free JSON API for Lobste.rs community news. No API key required for read access.

## Base URL

```
https://lobste.rs
```

## Authentication

None required for read-only access.

## Rate Limits

| Tier | Limit | Note |
|------|-------|------|
| Unauthenticated | ~60 req/min | Be polite. No documented hard cap |

**Etiquette**:
- Cache aggressively — stories change infrequently
- Use conditional requests if possible
- Avoid polling more than once per minute
- Set descriptive `User-Agent` header

## Endpoints

### Story Lists

| Endpoint | Description |
|----------|-------------|
| `GET /.json` | Hot stories (default front page) |
| `GET /newest.json` | Newest stories |
| `GET /hot.json` | Same as `/.json` — hot stories |
| `GET /s/{tag}.json` | Stories filtered by tag (e.g. `/s/rust.json`, `/s/programming.json`) |
| `GET /domain/{domain}.json` | Stories from a specific domain |
| `GET /page/{page}.json` | Pagination — page through hot stories |

### Individual Stories

| Endpoint | Description |
|----------|-------------|
| `GET /t/{short_id}.json` | Single story with comments |
| `GET /t/{short_id}/comments.json` | Flat comments for a story |

### Users

| Endpoint | Description |
|----------|-------------|
| `GET /u/{username}.json` | User profile + recent stories & comments |
| `GET /u/{username}/stories.json` | User's submitted stories |
| `GET /u/{username}/comments.json` | User's comments |

### Tags

| Endpoint | Description |
|----------|-------------|
| `GET /tags.json` | All available tags |

## Pagination

HATEOAS-style. Response includes `links` block with `next` URL:

```json
{
  "links": [
    { "href": "/newest.json?page=2", "rel": "next" }
  ]
}
```

Follow `links[0].href` to get next page. No manual page param needed — just follow the link.

## Response Format

### Story

```json
{
  "short_id": "abcdef",
  "short_id_url": "https://lobste.rs/s/abcdef",
  "created_at": "2024-01-15T10:30:00.000Z",
  "title": "Story Title",
  "url": "https://example.com/article",
  "score": 42,
  "comment_count": 7,
  "description": "Story description or text...",
  "comments_url": "https://lobste.rs/s/abcdef/comments",
  "submitter_user": {
    "username": "author",
    "created_at": "2023-01-01T00:00:00.000Z",
    "is_admin": false
  },
  "tags": ["programming", "rust"],
  "comments": []
}
```

| Field | Type | Description |
|-------|------|-------------|
| `short_id` | string | Unique story ID |
| `short_id_url` | string | Permalink |
| `created_at` | string | ISO 8601 timestamp |
| `title` | string | Story title |
| `url` | string | External URL or Lobste.rs if text |
| `score` | int | Upvote count |
| `comment_count` | int | Number of comments |
| `description` | string | Story text/markdown |
| `submitter_user` | object | `{username, created_at, is_admin}` |
| `tags` | [string] | Array of tags |
| `comments` | [Comment] | Full comment tree (when fetching single story) |
| `comment_count` | int | Total comment count |

### Comment

```json
{
  "short_id": "cdef01",
  "short_id_url": "https://lobste.rs/s/abcdef/cdef01",
  "created_at": "2024-01-15T11:00:00.000Z",
  "updated_at": "2024-01-15T11:05:00.000Z",
  "is_deleted": false,
  "is_moderated": false,
  "score": 5,
  "comment": "Comment body in markdown",
  "indent_level": 0,
  "hat": null,
  "issuer_id": null,
  "voting": { "upvoted": false, "downvoted": false }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `short_id` | string | Comment ID |
| `comment` | string | Body text (markdown) |
| `score` | int | Upvote count |
| `created_at` | string | ISO 8601 |
| `indent_level` | int | Nesting depth (0 = top-level) |
| `is_deleted` | bool | Deleted |
| `is_moderated` | bool | Moderated |
| `hat` | string | Hat/role badge, null if none |
| `voting` | object | `{upvoted, downvoted}` |

### Tag

```json
{
  "tag": "programming",
  "description": "General programming topics",
  "hotness_mod": 1.0,
  "is_media": false
}
```

## Example Usage

```bash
# Hot stories
curl "https://lobste.rs/hot.json"

# Newest stories
curl "https://lobste.rs/newest.json"

# Stories with tag
curl "https://lobste.rs/s/programming.json"

# Single story with comments
curl "https://lobste.rs/t/abcdef.json"

# User profile
curl "https://lobste.rs/u/username.json"

# All tags
curl "https://lobste.rs/tags.json"

# Paginate
curl "https://lobste.rs/newest.json?page=2"
```

## Limitations

- No full-text search API
- Max ~50 items per page
- Historical data limited — no archive beyond recent pages
- No API for `created_at` filtering or time-range queries
