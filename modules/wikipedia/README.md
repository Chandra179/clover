# Wikipedia Data Collector

Free JSON API over the MediaWiki Action API. No API key or authentication required for read access.

## Base URL

```
https://en.wikipedia.org/w/api.php
```

## Authentication

None required for read-only access. For editing operations (not covered here), OAuth or bot passwords are used.

## Format

Pass `format=json` on every request. Other formats: `xml`, `php`, `jsonfm` (pretty-print). Default without `format` parameter is `jsonfm` — always specify `format=json` for machine consumption.

## Rate Limits

| Tier | Limit |
|------|-------|
| Anonymous / normal user | 500 results per query (50 for slow queries); 50 multivalue inputs |
| Bots & admins (`apihighlimits`) | 5000 results per query (500 for slow queries); 500 inputs |

**Etiquette**: No hard rate limit, but:
- Don't make requests faster than 1–2 per second concurrently
- Use `maxlag` parameter on replicated-cluster wikis (Wikipedia) to avoid causing replication lag
- Set a descriptive `User-Agent` header (e.g. `CloverBot/0.1 (contact@example.com)`)
- Use `continue` for pagination, never page by offset

## Core Actions

### `action=query` — Main data fetcher

The most versatile endpoint. Fetches data via `prop`, `list`, and `meta` submodules.

```
GET /w/api.php?action=query&prop=extracts|info|pageimages&titles=Go_(programming_language)&format=json
```

#### Page Selection (mutually exclusive, pick one)

| Param | Type | Description |
|-------|------|-------------|
| `titles` | string | Pipe-separated page titles (max 50, 500 for bots) |
| `pageids` | int list | Pipe-separated page IDs |
| `revids` | int list | Revision IDs |
| `generator` | string | Query module to generate the page set (e.g. `search`, `random`, `categorymembers`) |

Generator param names must be prefixed with `g`. Example: `generator=search&gsrsearch=golang`.

**Titles with special characters**: URL-encode. Pipe character in title → use `%1F` (Unit Separator) as separator and prefix value with `%1F`.

#### `prop` — Page Properties

Get data about specific pages. Combine multiple using `|`.

| Prop | Key Params | Description |
|------|-----------|-------------|
| `extracts` | `exintro` (bool), `explaintext` (bool), `exchars` (int), `exlimit` (int, max 20) | Page text extracts. `exintro` = first section only. `explaintext` = plain text (no HTML) |
| `info` | `inprop=url|displaytitle|watched` | Basic page metadata |
| `revisions` | `rvprop=ids|timestamp|user|comment|content`, `rvlimit` (max 500), `rvdir=newer|older` | Revision history or page content. Use `rvslots=main` for current content |
| `pageimages` | `piprop=thumbnail|name|original`, `pithumbsize` (int), `pilimit` (max 50) | Page images / thumbnails |
| `categories` | `cllimit` (max 500), `clshow=hidden|!hidden` | Categories the page belongs to |
| `links` | `pllimit` (max 500), `plnamespace` (int) | Internal wiki links from the page |
| `langlinks` | `lllimit` (max 500), `lllang` (lang code) | Inter-language links |
| `templates` | `tllimit` (max 500) | Templates used on the page |
| `extlinks` | `ellimit` (max 500) | External URLs from the page |
| `images` | `imlimit` (max 500) | Images/files on the page |
| `coordinates` | `colimit` (max 500) | Geo coordinates |
| `contributors` | `pclimit` (max 500) | Page contributors |
| `pageviews` | `pvipdays` (int) | Daily pageview counts |
| `pageprops` | `ppprop` (string) | Page-level properties (disambiguation, wikibase_item, etc.) |
| `redirects` | — | Returns redirects to requested pages |

#### `list` — Lists (independent of page selection)

| List | Key Params | Description |
|------|-----------|-------------|
| `search` | `srsearch` (query), `srlimit` (max 500), `srwhat=text|title|nearmatch`, `srnamespace` | Full-text search |
| `categorymembers` | `cmtitle=Category:{name}`, `cmlimit` (max 500), `cmtype=page|subcat|file` | Pages in a category |
| `random` | `rnlimit` (max 20), `rnnamespace` (int) | Random pages |
| `allpages` | `apprefix` (prefix), `aplimit` (max 500), `apnamespace` | All pages alphabetically |
| `recentchanges` | `rclimit` (max 500), `rctype=edit|new|log` | Recent changes |
| `usercontribs` | `ucuser` (username), `uclimit` (max 500) | User's edits |
| `allcategories` | `aclimit`, `acprefix` | All categories |
| `allusers` | `aulimit`, `auprefix` | All users |
| `allimages` | `ailimit`, `aiprefix` | All files |
| `geosearch` | `gscoord=lat|lng`, `gsradius` (m), `gslimit` | Pages near coordinates |
| `prefixsearch` | `pssearch`, `pslimit` (max 500) | Title prefix search |
| `embeddedin` | `eititle`, `eilimit` | Pages that transclude a template |
| `backlinks` | `bltitle`, `bllimit` | Pages that link to a page |
| `logevents` | `letype`, `lelimit` (max 500) | Log entries |

#### `meta` — Site Metadata

| Meta | Key Params | Description |
|------|-----------|-------------|
| `siteinfo` | `siprop=general|namespaces|statistics|languages|extensions` | Wiki configuration |
| `userinfo` | `uiprop=groups|rights|email` | Current user info |
| `tokens` | `type=csrf|login|patrol|watch` | CSRF tokens for write ops |

### `action=parse` — Wikitext Parser

Parse wikitext to HTML.

```
GET /w/api.php?action=parse&page=Go_(programming_language)&prop=text|wikitext&format=json
```

| Param | Description |
|-------|-------------|
| `page` | Page title to parse |
| `pageid` | Page ID to parse |
| `oldid` | Revision ID to parse |
| `text` | Raw wikitext to parse (requires `title` or `contentmodel`) |
| `prop` | Properties: `text`, `wikitext`, `categories`, `links`, `templates`, `images`, `externallinks`, `sections`, `tocdata`, `langlinks`, `revid`, `displaytitle`, `headhtml`, `modules` |
| `section` | Parse only a specific section number or `new` |
| `disablelimitreport` | bool — Omit parser limit report |
| `disabletoc` | bool — Omit table of contents |
| `mobileformat` | bool — Mobile-friendly output |
| `useskin` | Skin: `vector`, `vector-2022`, `monobook`, `minerva`, `timeless` |

### `action=opensearch` — Quick Autocomplete

```
GET /w/api.php?action=opensearch&search=golang&limit=10&format=json
```

| Param | Description |
|-------|-------------|
| `search` | Search string |
| `limit` | Max results (default 10) |
| `namespace` | Namespace filter |

Response: `[query, [title1, title2, ...], [desc1, desc2, ...], [url1, url2, ...]]`

## Pagination (`continue`)

Most list/prop modules return a `continue` block for pagination:

```json
{
  "batchcomplete": "",
  "continue": {
    "sroffset": 10,
    "continue": "-||"
  },
  "query": { ... }
}
```

To get next page, merge `continue` params into your next request. Never page by offset.

## Response Format

```json
{
  "batchcomplete": "",
  "continue": { ... },
  "query": {
    "normalized": [
      { "from": "input_title", "to": "Normalized Title" }
    ],
    "pages": {
      "25039021": {
        "pageid": 25039021,
        "ns": 0,
        "title": "Go (programming language)",
        "extract": "Go is a high-level...",
        // ... more props
      }
    }
  }
}
```

### Key Fields

| Field | Description |
|-------|-------------|
| `batchcomplete` | Empty string if all results returned, absent if more available |
| `continue` | Present when more results exist; merge into next request |
| `query.normalized` | Title normalization info |
| `query.pages` | Map of pageid → page data |
| `query.searchinfo` | Search metadata (`totalhits`, suggestion) |
| `query.search` | Array of search results |
| `query.random` | Array of random page results |

### Page object fields

| Field | Type | Description |
|-------|------|-------------|
| `pageid` | int | Unique page ID |
| `ns` | int | Namespace ID (0 = main/article, 14 = category, etc.) |
| `title` | string | Page title |
| `extract` | string | Text extract (when `prop=extracts`) |
| `revisions` | array | Revision data (when `prop=revisions`) |
| `categories` | array | Category list (when `prop=categories`) |
| `links` | array | Internal links (when `prop=links`) |
| `thumbnail` | object | `{source, width, height}` (when `prop=pageimages`) |
| `pageimage` | string | Free image filename (when `prop=pageimages`) |
| `timestamp` | string | Last modified |
| `touched` | string | Last cache-touched timestamp |
| `length` | int | Page size in bytes |
| `redirect` | bool | True if this is a redirect (use `&redirects=` on query to resolve) |
| `contentmodel` | string | Content model (e.g. `wikitext`) |
| `pagelanguage` | string | Page language code |
| `fullurl` | string | Full page URL (when `inprop=url`) |
| `canonicalurl` | string | Canonical URL (when `inprop=url`) |

### Search result fields (`list=search`)

| Field | Type | Description |
|-------|------|-------------|
| `ns` | int | Namespace |
| `title` | string | Page title |
| `pageid` | int | Page ID |
| `size` | int | Page size in bytes |
| `wordcount` | int | Word count |
| `snippet` | string | HTML snippet with `<span class="searchmatch">` highlights |
| `timestamp` | string | Last modified timestamp |

## Namespaces (for `ns` field)

| ID | Name | Description |
|----|------|-------------|
| 0 | (Main) | Articles |
| 2 | User | User pages |
| 4 | Wikipedia | Project pages |
| 6 | File | File pages |
| 8 | MediaWiki | Interface messages |
| 10 | Template | Templates |
| 12 | Help | Help pages |
| 14 | Category | Categories |
| 828 | Module | Scribunto modules |

## CORS Support

For browser-based access, use `origin=*` parameter for anonymous requests. For authenticated CORS, use `origin=<domain>`.

## Error Handling

Errors return HTTP 200 with error body:

```json
{
  "error": {
    "code": "param-missing",
    "info": "Missing required parameter: titles",
    "*": "Messages..."
  }
}
```

Common error codes:
- `nosuchpageid` — Invalid page ID
- `missingtitle` — Title not found
- `ratelimited` — Too many requests
- `maxlag` — Database replication lag too high (wait and retry)

## Example Usage

```bash
# Search Wikipedia
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch=Monad+(functional+programming)&format=json&srlimit=5"

# Get page extract (intro only, plain text)
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=query&prop=extracts&exintro&explaintext&titles=Go_(programming_language)&format=json"

# Get page extract + info + images
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=query&prop=extracts|info|pageimages&exintro&explaintext&inprop=url&pithumbsize=200&titles=Go_(programming_language)&format=json"

# Get full page content (via revisions)
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&rvslots=main&titles=Go_(programming_language)&format=json"

# Parse page to HTML
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=parse&page=Go_(programming_language)&prop=text&format=json"

# Get pages in a category
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=query&list=categorymembers&cmtitle=Category:Programming_languages&cmlimit=50&format=json"

# Get random articles
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=query&list=random&rnlimit=5&rnnamespace=0&format=json"

# Get pageview stats
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=query&prop=pageviews&titles=Go_(programming_language)&format=json"

# OpenSearch autocomplete
curl -H "User-Agent: Clover/0.1" \
  "https://en.wikipedia.org/w/api.php?action=opensearch&search=golang&limit=5&format=json"
```

## Limitations

- Max 50 pages per `titles`/`pageids` request (500 for bots)
- Max 500 results per `prop`/`list` query (5000 for bots, 50 for "slow" queries like `search` → 50 max)
- Page extracts `exlimit` max is 20
- Search via `list=search` uses `srlimit` max 50 (500 for bots)
- No means to bulk-download all pages at once (use [Wikimedia dumps](https://dumps.wikimedia.org/) for that)
- `prop=extracts` returns text up to ~12000 chars (or specified by `exchars`)

## Wikipedia Alternatives

| Alternative | URL | Auth | Rate |
|-------------|-----|------|------|
| MediaWiki API (any wiki) | `{wiki}/w/api.php` | None | Same as Wikipedia |
| Wikimedia REST API | `https://en.wikipedia.org/api/rest_v1/` | None | 200 req/s |
| Wikimedia Enterprise API | `https://api.enterprise.wikimedia.com/` | Key required | Structured/snapshot access |
