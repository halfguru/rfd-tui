package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfdtui/internal/types"
)

const defaultBaseURL = "https://forums.redflagdeals.com"

const (
	maxRetries = 3
	baseDelay  = 500 * time.Millisecond
	maxDelay   = 5 * time.Second
)

type Client struct {
	http      *http.Client
	baseURL   string
	userCache map[int]string
	mu        sync.RWMutex
}

func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL:   defaultBaseURL,
		userCache: make(map[int]string),
	}
}

func NewWithBaseURL(url string) *Client {
	c := New()
	c.baseURL = url
	return c
}

type UsernameMsg struct {
	AuthorID int
	Username string
}

type TopicsMsg struct {
	Topics []types.Topic
	Page   int
	Pager  types.Pager
}

type PostsMsg struct {
	Posts []types.Post
	Page  int
	Pager types.Pager
}

type ErrMsg struct {
	Err error
}

type BrowserOpenMsg struct {
	URL string
	Err error
}

func (c *Client) FetchTopics(page int) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("%s/api/topics?forum_id=9&per_page=40&page=%d", c.baseURL, page)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return c.doGetWithRetry(ctx, url, func(body json.Decoder) tea.Msg {
			var result types.TopicsResponse
			if err := body.Decode(&result); err != nil {
				return ErrMsg{Err: fmt.Errorf("failed to parse response: %w", err)}
			}
			return TopicsMsg{Topics: result.Topics, Page: result.Pager.Page, Pager: result.Pager}
		})
	}
}

func (c *Client) FetchPosts(topicID, page int) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("%s/api/topics/%d/posts?per_page=40&page=%d", c.baseURL, topicID, page)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return c.doGetWithRetry(ctx, url, func(body json.Decoder) tea.Msg {
			var result types.PostsResponse
			if err := body.Decode(&result); err != nil {
				return ErrMsg{Err: fmt.Errorf("failed to parse response: %w", err)}
			}
			return PostsMsg{Posts: result.Posts, Page: result.Pager.Page, Pager: result.Pager}
		})
	}
}

func (c *Client) FetchUsername(authorID int) tea.Cmd {
	return func() tea.Msg {
		c.mu.RLock()
		name, ok := c.userCache[authorID]
		c.mu.RUnlock()
		if ok {
			return UsernameMsg{AuthorID: authorID, Username: name}
		}
		url := fmt.Sprintf("%s/api/users/%d", c.baseURL, authorID)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return c.doGetWithRetry(ctx, url, func(body json.Decoder) tea.Msg {
			var result types.UserResponse
			if err := body.Decode(&result); err != nil {
				return UsernameMsg{AuthorID: authorID, Username: fmt.Sprintf("user_%d", authorID)}
			}
			c.mu.Lock()
			c.userCache[authorID] = result.User.Username
			c.mu.Unlock()
			return UsernameMsg{AuthorID: authorID, Username: result.User.Username}
		})
	}
}

func (c *Client) CachedUsername(authorID int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if name, ok := c.userCache[authorID]; ok {
		return name
	}
	return fmt.Sprintf("user_%d", authorID)
}

func (c *Client) doGetWithRetry(ctx context.Context, url string, handler func(json.Decoder) tea.Msg) tea.Msg {
	var lastErr error
	for attempt := range maxRetries {
		if attempt > 0 {
			delay := backoff(attempt)
			select {
			case <-ctx.Done():
				return ErrMsg{Err: ctx.Err()}
			case <-time.After(delay):
			}
		}

		msg := c.doGet(ctx, url, handler)
		if _, ok := msg.(ErrMsg); !ok {
			return msg
		}
		lastErr = msg.(ErrMsg).Err
	}
	return ErrMsg{Err: fmt.Errorf("after %d retries: %w", maxRetries, lastErr)}
}

func (c *Client) doGet(ctx context.Context, url string, handler func(json.Decoder) tea.Msg) tea.Msg {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return ErrMsg{Err: err}
	}
	req.Header.Set("User-Agent", "rfdtui/1.0 (terminal deal browser)")

	resp, err := c.http.Do(req)
	if err != nil {
		return ErrMsg{Err: err}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests {
		return ErrMsg{Err: fmt.Errorf("rate limited (429)")}
	}
	if resp.StatusCode >= 500 {
		return ErrMsg{Err: fmt.Errorf("server error: %d", resp.StatusCode)}
	}
	if resp.StatusCode != http.StatusOK {
		return ErrMsg{Err: fmt.Errorf("API returned status %d", resp.StatusCode)}
	}

	return handler(*json.NewDecoder(resp.Body))
}

func backoff(attempt int) time.Duration {
	delay := baseDelay * time.Duration(1<<uint(attempt))
	delay = time.Duration(float64(delay) * (0.8 + rand.Float64()*0.4))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}
