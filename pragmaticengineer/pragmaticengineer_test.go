package pragmaticengineer_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/pragmaticengineer-cli/pragmaticengineer"
)

const fakeJSON = `[
  {
    "title": "Test Post One",
    "subtitle": "A subtitle",
    "post_date": "2026-06-11T16:26:58.630Z",
    "canonical_url": "https://example.com/p/test-post-one",
    "audience": "only_paid",
    "slug": "test-post-one"
  },
  {
    "title": "Test Post Two",
    "subtitle": "Another subtitle",
    "post_date": "2026-06-01T12:00:00.000Z",
    "canonical_url": "https://example.com/p/test-post-two",
    "audience": "everyone",
    "slug": "test-post-two"
  }
]`

func newTestClient(ts *httptest.Server) *pragmaticengineer.Client {
	cfg := pragmaticengineer.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return pragmaticengineer.NewClient(cfg)
}

func TestTop(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fakeJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	posts, err := c.Top(context.Background(), 25, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 {
		t.Fatalf("want 2 posts, got %d", len(posts))
	}
	if posts[0].Title != "Test Post One" {
		t.Errorf("Title = %q", posts[0].Title)
	}
	if posts[0].Date != "2026-06-11" {
		t.Errorf("Date = %q, want 2026-06-11", posts[0].Date)
	}
	if posts[0].Audience != "paid" {
		t.Errorf("Audience = %q, want paid", posts[0].Audience)
	}
	if posts[1].Audience != "free" {
		t.Errorf("Audience = %q, want free", posts[1].Audience)
	}
	if posts[0].Rank != 1 {
		t.Errorf("Rank = %d, want 1", posts[0].Rank)
	}
}
