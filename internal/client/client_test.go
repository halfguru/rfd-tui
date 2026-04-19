package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/halfguru/rfd-tui/internal/types"
)

func TestFetchTopics(t *testing.T) {
	resp := types.TopicsResponse{
		Pager: types.Pager{Page: 1, TotalPages: 2, Total: 80},
		Topics: []types.Topic{
			{TopicID: 1, Title: "Test Deal", Score: 42},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("forum_id") != "9" {
			t.Errorf("expected forum_id=9, got %s", r.URL.Query().Get("forum_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatal(err)
		}
	}))
	defer srv.Close()

	c := NewWithBaseURL(srv.URL)
	cmd := c.FetchTopics(1)
	msg := cmd()

	topicsMsg, ok := msg.(TopicsMsg)
	if !ok {
		t.Fatalf("expected TopicsMsg, got %T", msg)
	}
	if len(topicsMsg.Topics) != 1 {
		t.Errorf("expected 1 topic, got %d", len(topicsMsg.Topics))
	}
	if topicsMsg.Topics[0].Title != "Test Deal" {
		t.Errorf("expected 'Test Deal', got %q", topicsMsg.Topics[0].Title)
	}
}

func TestFetchPosts(t *testing.T) {
	resp := types.PostsResponse{
		Pager: types.Pager{Page: 1, TotalPages: 1},
		Posts: []types.Post{
			{PostID: 1, Body: "<p>Hello</p>", AuthorID: 10},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatal(err)
		}
	}))
	defer srv.Close()

	c := NewWithBaseURL(srv.URL)
	cmd := c.FetchPosts(123, 1)
	msg := cmd()

	postsMsg, ok := msg.(PostsMsg)
	if !ok {
		t.Fatalf("expected PostsMsg, got %T", msg)
	}
	if len(postsMsg.Posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(postsMsg.Posts))
	}
}

func TestRetryOnServerError(t *testing.T) {
	var attempts atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := types.TopicsResponse{
			Pager:  types.Pager{Page: 1},
			Topics: []types.Topic{{TopicID: 1}},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatal(err)
		}
	}))
	defer srv.Close()

	c := NewWithBaseURL(srv.URL)
	cmd := c.FetchTopics(1)
	msg := cmd()

	topicsMsg, ok := msg.(TopicsMsg)
	if !ok {
		t.Fatalf("expected TopicsMsg after retry, got %T", msg)
	}
	if len(topicsMsg.Topics) != 1 {
		t.Errorf("expected 1 topic after retry, got %d", len(topicsMsg.Topics))
	}
	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestCachedUsername(t *testing.T) {
	c := New()

	name := c.CachedUsername(42)
	if name != "user_42" {
		t.Errorf("expected 'user_42', got %q", name)
	}

	c.mu.Lock()
	c.userCache[42] = "testuser"
	c.mu.Unlock()

	name = c.CachedUsername(42)
	if name != "testuser" {
		t.Errorf("expected 'testuser', got %q", name)
	}
}

func TestBackoff(t *testing.T) {
	for attempt := 1; attempt <= 5; attempt++ {
		delay := backoff(attempt)
		if delay < baseDelay {
			t.Errorf("backoff(%d) = %v, want >= %v", attempt, delay, baseDelay)
		}
		if delay > maxDelay {
			t.Errorf("backoff(%d) = %v, want <= %v", attempt, delay, maxDelay)
		}
	}
}
