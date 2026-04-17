package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/simon/rfd/internal/types"
	tea "charm.land/bubbletea/v2"
)

const baseURL = "https://forums.redflagdeals.com"

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type TopicsMsg struct {
	Topics []types.Topic
	Page   int
	Pager  types.Pager
}

type ErrMsg struct {
	Err error
}

func (c *Client) FetchTopics(page int) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("%s/api/topics?forum_id=9&per_page=40&page=%d", baseURL, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ErrMsg{Err: err}
		}
		req.Header.Set("User-Agent", "rfd-tui/1.0 (terminal deal browser)")

		resp, err := c.http.Do(req)
		if err != nil {
			return ErrMsg{Err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return ErrMsg{Err: fmt.Errorf("API returned status %d", resp.StatusCode)}
		}

		var result types.TopicsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return ErrMsg{Err: fmt.Errorf("failed to parse response: %w", err)}
		}

		return TopicsMsg{
			Topics: result.Topics,
			Page:   result.Pager.Page,
			Pager:  result.Pager,
		}
	}
}
