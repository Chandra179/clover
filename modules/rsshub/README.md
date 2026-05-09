# RSSHub Data Collector

Proxy for any RSSHub instance. RSSHub converts RSS feeds from 1000+ sources into JSON. Self-host an instance or use public one. No API key required.

## Base URL

```
https://rsshub.app  (public instance)
```

Or self-hosted instance. Configurable per deployment.

## Authentication

None required. RSSHub instances are public by default. Some routes may require cookies for sites behind auth walls.

## Rate Limits

| Instance | Limit | Note |
|----------|-------|------|
| `rsshub.app` | ~250 req/3min per IP | Public instance, enforced by Cloudflare |
| Self-hosted | Unlimited | Controlled by your own infrastructure |

**Etiquette**:
- Prefer self-hosting for production use
- Cache responses aggressively (RSSHub is slow — it scrapes live sites)
- rsshub.app is a free community service — don't abuse it

## Route Format

```
GET /{route}.json
```

Route path mirrors the RSSHub route table: `/source/category/param1/param2.json`

### Common Routes

| Route | Description |
|-------|-------------|
| `GET /hackernews.json` | HackerNews front page |
| `GET /hackernews/new.json` | HackerNews new |
| `GET /hackernews/best.json` | HackerNews best |
| `GET /github/trending/daily.json` | GitHub trending repos (daily) |
| `GET /github/trending/weekly.json` | GitHub trending (weekly) |
| `GET /github/trending/monthly.json` | GitHub trending (monthly) |
| `GET /github/trending/daily/rust.json` | GitHub trending by language |
| `GET /programming/dev.to.json` | dev.to articles |
| `GET /zhihu/pin/daily.json` | Zhihu daily pins |
| `GET /v2ex/topics/latest.json` | V2EX latest |
| `GET /solidot/linux.json` | Solidot Linux news |
| `GET /36kr/motif/{id}.json` | 36Kr motif articles |
| `GET /nytimes/daily.json` | NYT daily |
| `GET /bbc/world.json` | BBC World |
| `GET /reuters/topics/{topic}.json` | Reuters by topic |

Full route list: [https://docs.rsshub.app/routes](https://docs.rsshub.app/routes)

## Response Format

```json
{
  "title": "GitHub Trending",
  "link": "https://github.com/trending",
  "description": "GitHub Trending Daily",
  "lastBuildDate": "Mon, 15 Jan 2024 10:00:00 GMT",
  "items": [
    {
      "title": "owner/repo — description",
      "link": "https://github.com/owner/repo",
      "guid": "https://github.com/owner/repo",
      "description": "<p>...</p>",
      "author": "owner",
      "pubDate": "Mon, 15 Jan 2024 10:00:00 GMT"
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Feed title |
| `link` | string | Source URL |
| `description` | string | Feed description |
| `lastBuildDate` | string | RFC 1123 timestamp |
| `items` | [Item] | Array of feed items |

### Item Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Article title |
| `link` | string | Article URL |
| `guid` | string | Unique ID |
| `description` | string | HTML content |
| `author` | string | Author name |
| `pubDate` | string | RFC 1123 timestamp |
| `category` | string | Article category/tag |

## Query Parameters

| Param | Description |
|-------|-------------|
| `?limit=N` | Limit results (varies by route, default ~20-50) |
| `?lang={code}` | Language filter (route-dependent) |

## Caching

RSSHub caches aggressively:
- Default TTL: 5 minutes (configurable on self-hosted instances)
- Cache key = full URL path
- Cache stored in memory by default (Redis option available)

## Example Usage

```bash
# Public instance — HackerNews
curl "https://rsshub.app/hackernews.json"

# Public instance — GitHub trending Rust
curl "https://rsshub.app/github/trending/daily/rust.json"

# Public instance — BBC World
curl "https://rsshub.app/bbc/world.json"

# Self-hosted instance
curl "http://localhost:1200/hackernews.json"
```

## Limitations

- Response time depends on upstream site scrape speed (can be 2-10s)
- Some routes require cookies or special headers
- Public rsshub.app may be unreliable under load
- No pagination — RSSHub returns whatever RSS limited to
- Some sources have broken/incomplete descriptions
- Not all routes available on all instances (depends on RSSHub version and config)
