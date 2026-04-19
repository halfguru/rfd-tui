package views

import (
	"strings"
	"testing"
	"time"
)

func TestRelativeAge(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{"just now", time.Now(), "just now"},
		{"30s ago", time.Now().Add(-30 * time.Second), "just now"},
		{"5m ago", time.Now().Add(-5 * time.Minute), "5m ago"},
		{"1h ago", time.Now().Add(-1 * time.Hour), "1h ago"},
		{"23h ago", time.Now().Add(-23 * time.Hour), "23h ago"},
		{"1d ago", time.Now().Add(-25 * time.Hour), "1d ago"},
		{"7d ago", time.Now().Add(-7 * 24 * time.Hour), "7d ago"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relativeAge(tt.time)
			if got != tt.want {
				t.Errorf("relativeAge() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderScrollBar(t *testing.T) {
	tests := []struct {
		name   string
		cursor int
		total  int
		height int
		empty  bool
	}{
		{"empty", 0, 0, 20, true},
		{"single item", 0, 1, 20, false},
		{"mid position", 5, 10, 20, false},
		{"last item", 9, 10, 20, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderScrollBar(tt.cursor, tt.total, tt.height)
			if tt.empty {
				if got != "" {
					t.Errorf("expected empty, got %q", got)
				}
			} else {
				if !strings.Contains(got, "▌") {
					t.Errorf("expected thumb in scrollbar, got %q", got)
				}
			}
		})
	}
}

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain text", "hello", "hello"},
		{"bold", "<b>hello</b>", "hello"},
		{"br tag", "line1<br>line2", "line1\nline2"},
		{"br self-close", "line1<br/>line2", "line1\nline2"},
		{"entities", "&amp;&quot;&nbsp;", "&\""},
		{"nested", "<div><p>text</p></div>", "text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripHTML(tt.input)
			if got != tt.want {
				t.Errorf("stripHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  string
	}{
		{"short line", "hello", 80, "hello\n"},
		{"exact width", "hello", 5, "hello\n"},
		{"wrap", "hello world foo", 8, "hello\nworld\nfoo\n"},
		{"zero width uses default", "test", 0, "test\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapText(tt.input, tt.width)
			if got != tt.want {
				t.Errorf("wrapText() = %q, want %q", got, tt.want)
			}
		})
	}
}
