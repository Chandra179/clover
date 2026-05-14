package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"brook/modules/common"

	"github.com/gin-gonic/gin"
)

type mockFetcher struct {
	name    string
	results []common.CategoryResult
	err     error
}

func (m *mockFetcher) FetchCategory(category string) ([]common.CategoryResult, error) {
	return m.results, m.err
}

func setupTest(fetchers []common.Fetcher) *gin.Engine {
	gin.SetMode(gin.TestMode)
	d := NewDependencies(fetchers)
	r := gin.New()
	r.GET("/news", d.NewsHandler)
	return r
}

func TestNewsHandler_Success(t *testing.T) {
	fetchers := []common.Fetcher{
		&mockFetcher{
			name: "source1",
			results: []common.CategoryResult{
				{Title: "Article 1", URL: "http://example.com/1", Category: "tech", Source: "source1"},
				{Title: "Article 2", URL: "http://example.com/2", Category: "tech", Source: "source1"},
			},
		},
		&mockFetcher{
			name: "source2",
			results: []common.CategoryResult{
				{Title: "Article 3", URL: "http://example.com/3", Category: "tech", Source: "source2"},
			},
		},
	}

	r := setupTest(fetchers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news?category=tech", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp NewsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if resp.Count != 3 {
		t.Errorf("expected 3 items, got %d", resp.Count)
	}
	if resp.Category != "tech" {
		t.Errorf("expected category=tech, got %s", resp.Category)
	}
	if len(resp.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(resp.Sources))
	}
}

func TestNewsHandler_DefaultCategory(t *testing.T) {
	fetchers := []common.Fetcher{
		&mockFetcher{
			name: "source1",
			results: []common.CategoryResult{
				{Title: "Article 1", URL: "http://example.com/1", Category: "tech", Source: "source1"},
			},
		},
	}

	r := setupTest(fetchers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp NewsResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Category != "tech" {
		t.Errorf("expected default category=tech, got %s", resp.Category)
	}
}

func TestNewsHandler_PartialError(t *testing.T) {
	fetchers := []common.Fetcher{
		&mockFetcher{
			name: "good",
			results: []common.CategoryResult{
				{Title: "Article 1", URL: "http://example.com/1", Category: "tech", Source: "good"},
			},
		},
		&mockFetcher{
			name: "bad",
			err:  http.ErrAbortHandler,
		},
	}

	r := setupTest(fetchers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news?category=tech", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (partial ok)", w.Code)
	}

	var resp NewsResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Count != 1 {
		t.Errorf("expected 1 item, got %d", resp.Count)
	}
}

func TestNewsHandler_AllErrors(t *testing.T) {
	fetchers := []common.Fetcher{
		&mockFetcher{name: "bad1", err: http.ErrAbortHandler},
		&mockFetcher{name: "bad2", err: http.ErrAbortHandler},
	}

	r := setupTest(fetchers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news?category=tech", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestNewsHandler_EmptyResults(t *testing.T) {
	fetchers := []common.Fetcher{
		&mockFetcher{
			name:    "empty",
			results: []common.CategoryResult{},
		},
	}

	r := setupTest(fetchers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news?category=tech", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp NewsResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Count != 0 {
		t.Errorf("expected 0 items, got %d", resp.Count)
	}
}
