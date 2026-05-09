# HackerNews Data Collector

Free JSON API via Firebase Realtime Database. No API key required for read access.

## Base URL

```
https://hacker-news.firebaseio.com/v0/
```

## Authentication

None required. Read-only Firebase REST API.

## Rate Limits

Undocumented by Firebase/HN. Known behavior:
- ~500 req/min per IP before throttling
- HTTP 429 possible under burst load
- Cache aggressively — stories change infrequently (new list updates every ~minute)
- Use `ETag`/`If-None-Match` if client supports it

## Endpoints

### Story Lists

| Endpoint | Description |
|----------|-------------|
| `GET /topstories.json` | Up to 500 top & trending story IDs |
| `GET /newstories.json` | Up to 500 newest story IDs |
| `GET /beststories.json` | Up to 500 highest-rated story IDs |
| `GET /askstories.json` | Up to 200 "Ask HN" story IDs |
| `GET /showstories.json` | Up to 200 "Show HN" story IDs |
| `GET /jobstories.json` | Up to 200 "Job HN" story IDs |

Returns: array of int IDs. Max 500 items. Order approximate.

### Individual Items

| Endpoint | Description |
|----------|-------------|
| `GET /item/{id}.json` | Story, comment, ask, show, job, or poll |

### Users

| Endpoint | Description |
|----------|-------------|
| `GET /user/{id}.json` | User profile + submitted stories & comments |

### Max Item ID

| Endpoint | Description |
|----------|-------------|
| `GET /maxitem.json` | Current largest item ID (monotonically increasing) |

### Changes & Profiles

| Endpoint | Description |
|----------|-------------|
| `GET /updates.json` | Recently changed items and users |
| `GET /newprompts.json` | Newest GPT prompts |

## Response Format

### Item (type field)

```json
{
  "id": 8863,
  "deleted": false,
  "type": "story",
  "by": "dhouston",
  "time": 1175714200,
  "text": "Ask HN: What would you ...",
  "dead": false,
  "parent": 0,
  "poll": 0,
  "kids": [8952, 9224, 8917, ...],
  "url": "http://example.com",
  "score": 104,
  "title": "My YC app: Dropbox",
  "parts": [1, 2, 3],
  "descendants": 71
}
```

| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Unique story/comment ID |
| `type` | string | `"story"`, `"comment"`, `"poll"`, `"pollopt"`, `"job"` |
| `by` | string | Author username |
| `time` | int | Unix timestamp |
| `text` | string | HTML body (comments, Ask HN text) |
| `url` | string | External link (stories) |
| `title` | string | Story title |
| `score` | int | Story points |
| `descendants` | int | Comment count |
| `kids` | [int] | Child comment IDs (top-level only) |
| `parent` | int | Parent comment ID |
| `parts` | [int] | Poll option IDs |
| `deleted` | bool | Deleted item |
| `dead` | bool | Flagged dead |

### User

```json
{
  "id": "jl",
  "delay": 0,
  "created": 1173923446,
  "karma": 2937,
  "about": "Description text",
  "submitted": [1, 2, 3]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Username |
| `created` | int | Account creation timestamp |
| `karma` | int | Karma score |
| `about` | string | Profile text |
| `submitted` | [int] | Array of submitted item IDs |
| `delay` | int | Minutes of submission delay |

### Story List

```
[ 39440273, 39439892, 39439761, ... ]
```

Max 500 IDs. Order = approximate rank.

### Updates

```json
{
  "items": [8423305, 8420805, ...],
  "profiles": ["thegreatergood", "b6a1234f", ...]
}
```

## Caching Notes

- `/topstories.json` changes every ~few minutes
- Individual `/item/{id}.json` is immutable after ~30s of creation
- Cache stale data for 30-60s on client side
- Firebase serves stale data from CDN edge — TTL ~30s

## Example Usage

```bash
# Top stories
curl "https://hacker-news.firebaseio.com/v0/topstories.json"

# Item details
curl "https://hacker-news.firebaseio.com/v0/item/39440273.json"

# User details
curl "https://hacker-news.firebaseio.com/v0/user/pg.json"

# Recent updates
curl "https://hacker-news.firebaseio.com/v0/updates.json"

# Max item ID
curl "https://hacker-news.firebaseio.com/v0/maxitem.json"
```

## Limitations

- Max 500 IDs per list endpoint — no pagination
- No search endpoint (use Algolia HN Search API `https://hn.algolia.com/api/v1/` for search)
- No filter by score/time beyond what lists provide
- Comment tree must be reconstructed via `kids` chain — no nested comment response
- No CORS support for browser JS
- Items sometimes missing for ~seconds after creation (Firebase eventual consistency)
