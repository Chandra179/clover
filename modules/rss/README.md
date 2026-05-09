# Generic RSS/Atom Feed Parser

Fetches and parses RSS 2.0, Atom, and RDF feeds. No API key required. Feed URLs configured by category.

## No Authentication

RSS/Atom feeds are public XML documents. No API key or auth required.

## Rate Limits

Depends entirely on upstream server. Be polite:
- Cache aggressively (min 5 min between fetches per feed)
- Respect `Retry-After` headers
- Set descriptive `User-Agent`

## Supported Formats

| Format | Detected By |
|--------|-------------|
| RSS 2.0 | `<rss version="2.0">` |
| Atom | `<feed xmlns="http://www.w3.org/2005/Atom">` |
| RDF/RSS 1.0 | `<rdf:RDF>` |

## Configuration

Each category maps to a list of feed URLs in config:

```yaml
rss:
  base_url: ""  # not used for this module
  categories:
    tech:
      - https://news.ycombinator.com/rss
      - https://lobste.rs/rss
      - https://blog.golang.org/feed.atom
    economy:
      - https://feeds.bbci.co.uk/news/business/rss.xml
```

## Response Format

### Standard RSS 2.0

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Feed Title</title>
    <link>https://example.com</link>
    <description>Feed description</description>
    <item>
      <title>Article Title</title>
      <link>https://example.com/article</link>
      <guid>unique-id</guid>
      <pubDate>Mon, 15 Jan 2024 10:00:00 GMT</pubDate>
      <description>Article description or content</description>
      <author>author@example.com (Name)</author>
      <category>Tag</category>
    </item>
  </channel>
</rss>
```

### Atom

```xml
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Feed Title</title>
  <link href="https://example.com"/>
  <entry>
    <title>Article Title</title>
    <link href="https://example.com/article"/>
    <id>unique-id</id>
    <published>2024-01-15T10:00:00Z</published>
    <updated>2024-01-15T10:00:00Z</updated>
    <summary>Description</summary>
    <author><name>Author</name></author>
    <category term="Tag"/>
  </entry>
</feed>
```

### Normalized Output

All formats normalized to:

```json
{
  "title": {
    "title": "Article Title",
    "url": "https://example.com/article",
    "content": "Article description",
    "category": "tech",
    "source": "rss"
  }
}
```

## Example Usage

```bash
# RSS 2.0
curl "https://news.ycombinator.com/rss"

# Atom
curl "https://blog.golang.org/feed.atom"

# RDF
curl "https://web.mit.edu/newsoffice/feed/rss"
```

## Limitations

- No server-side filtering — parse client-side
- RSS descriptions truncated if very long
- Some feeds serve only summaries, not full content
- Atom entries may contain HTML in `<summary>` or `<content>` elements
- Some feeds have malformed XML — parser is lenient but may fail
- No built-in deduplication across multiple feed URLs
