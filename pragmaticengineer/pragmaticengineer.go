// Package pragmaticengineer is the library behind the pragmaticengineer command line:
// the HTTP client, request shaping, and the typed data models for The Pragmatic Engineer.
//
// The Client here is the spine every command shares. It sets a real
// User-Agent, paces requests so a busy session stays polite, and retries the
// transient failures (429 and 5xx) that any public site throws under load.
package pragmaticengineer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Host is the site this client talks to.
const Host = "newsletter.pragmaticengineer.com"

// BaseURL is the root every request is built from.
const BaseURL = "https://" + Host

// Config holds tunables for the Client.
type Config struct {
	BaseURL   string
	Rate      time.Duration
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   BaseURL,
		Rate:      500 * time.Millisecond,
		Timeout:   30 * time.Second,
		Retries:   3,
		UserAgent: "Mozilla/5.0 (compatible; pragmaticengineer-cli/0.1)",
	}
}

// Client talks to The Pragmatic Engineer over HTTP.
type Client struct {
	cfg    Config
	http   *http.Client
	lastAt time.Time
}

// NewClient returns a Client configured from cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}


func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	if c.cfg.Rate > 0 {
		if wait := c.cfg.Rate - time.Since(c.lastAt); wait > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
		}
	}
	url := c.cfg.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	var (
		resp *http.Response
		body []byte
	)
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		resp, err = c.http.Do(req)
		if err != nil {
			if attempt < c.cfg.Retries {
				continue
			}
			return nil, err
		}
		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		c.lastAt = time.Now()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			if attempt < c.cfg.Retries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		return body, nil
	}
	return body, nil
}

type apiPost struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	PostDate string `json:"post_date"`
	CanonURL string `json:"canonical_url"`
	Audience string `json:"audience"`
	Slug     string `json:"slug"`
}

// Top fetches recent posts from The Pragmatic Engineer Substack API.
func (c *Client) Top(ctx context.Context, limit, offset int) ([]*Post, error) {
	if limit <= 0 {
		limit = 25
	}
	path := fmt.Sprintf("/api/v1/posts?limit=%d&offset=%d", limit, offset)
	body, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	var raw []apiPost
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	posts := make([]*Post, 0, len(raw))
	for i, p := range raw {
		date := p.PostDate
		if len(date) >= 10 {
			date = date[:10]
		}
		audience := p.Audience
		if audience == "everyone" {
			audience = "free"
		} else if strings.Contains(audience, "paid") {
			audience = "paid"
		}
		posts = append(posts, &Post{
			Rank:     offset + i + 1,
			Date:     date,
			Audience: audience,
			Title:    p.Title,
			Subtitle: p.Subtitle,
			URL:      p.CanonURL,
		})
	}
	return posts, nil
}
