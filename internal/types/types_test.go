package types

import (
	"testing"
)

func TestTopic_DealerName(t *testing.T) {
	tests := []struct {
		name  string
		topic Topic
		want  string
	}{
		{"nil offer", Topic{}, ""},
		{"with dealer", Topic{Offer: &Offer{DealerName: "Amazon"}}, "Amazon"},
		{"empty dealer", Topic{Offer: &Offer{DealerName: ""}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.topic.DealerName(); got != tt.want {
				t.Errorf("DealerName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTopic_DealURL(t *testing.T) {
	tests := []struct {
		name  string
		topic Topic
		want  string
	}{
		{"offer url", Topic{Offer: &Offer{URL: "https://amzn.com/x"}, WebPath: "/x"}, "https://amzn.com/x"},
		{"fallback url", Topic{WebPath: "/threads/123"}, "https://forums.redflagdeals.com/threads/123"},
		{"empty offer url", Topic{Offer: &Offer{URL: ""}, WebPath: "/threads/456"}, "https://forums.redflagdeals.com/threads/456"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.topic.DealURL(); got != tt.want {
				t.Errorf("DealURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTopic_Price(t *testing.T) {
	tests := []struct {
		name  string
		topic Topic
		want  string
	}{
		{"nil offer", Topic{}, ""},
		{"with price", Topic{Offer: &Offer{Price: "$99"}}, "$99"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.topic.Price(); got != tt.want {
				t.Errorf("Price() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTopic_Savings(t *testing.T) {
	tests := []struct {
		name  string
		topic Topic
		want  string
	}{
		{"nil offer", Topic{}, ""},
		{"with savings", Topic{Offer: &Offer{Savings: "50%"}}, "50%"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.topic.Savings(); got != tt.want {
				t.Errorf("Savings() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTopic_CategoryID(t *testing.T) {
	tests := []struct {
		name  string
		topic Topic
		want  int
	}{
		{"nil offer", Topic{}, 0},
		{"with category", Topic{Offer: &Offer{CategoryID: 9}}, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.topic.CategoryID(); got != tt.want {
				t.Errorf("CategoryID() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTopic_CategoryName(t *testing.T) {
	tests := []struct {
		name  string
		topic Topic
		want  string
	}{
		{"nil offer", Topic{}, ""},
		{"computers", Topic{Offer: &Offer{CategoryID: 9}}, "Computers & Electronics"},
		{"unknown", Topic{Offer: &Offer{CategoryID: 999}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.topic.CategoryName(); got != tt.want {
				t.Errorf("CategoryName() = %q, want %q", got, tt.want)
			}
		})
	}
}
