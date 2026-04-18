package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfd/internal/types"
)

const baseURL = "https://forums.redflagdeals.com"

type Client struct {
	http      *http.Client
	userCache map[int]string
	mu        sync.RWMutex
}

func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
		userCache: make(map[int]string),
	}
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
		url := fmt.Sprintf("%s/api/topics?forum_id=9&per_page=40&page=%d", baseURL, page)
		return c.doGet(url, func(body json.Decoder) tea.Msg {
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
		url := fmt.Sprintf("%s/api/topics/%d/posts?per_page=40&page=%d", baseURL, topicID, page)
		return c.doGet(url, func(body json.Decoder) tea.Msg {
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
		url := fmt.Sprintf("%s/api/users/%d", baseURL, authorID)
		return c.doGet(url, func(body json.Decoder) tea.Msg {
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

func (c *Client) doGet(url string, handler func(json.Decoder) tea.Msg) tea.Msg {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ErrMsg{Err: err}
	}
	req.Header.Set("User-Agent", "rfd-tui/1.0 (terminal deal browser)")

	resp, err := c.http.Do(req)
	if err != nil {
		return ErrMsg{Err: err}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return ErrMsg{Err: fmt.Errorf("API returned status %d", resp.StatusCode)}
	}

	return handler(*json.NewDecoder(resp.Body))
}
