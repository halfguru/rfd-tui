package types

import "time"

type TopicsResponse struct {
	Pager  Pager   `json:"pager"`
	Topics []Topic `json:"topics"`
}

type Pager struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	TotalOnPage int `json:"total_on_page"`
	TotalPages  int `json:"total_pages"`
	Page        int `json:"page"`
}

type Topic struct {
	TopicID      int       `json:"topic_id"`
	ForumID      int       `json:"forum_id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	PostTime     time.Time `json:"post_time"`
	LastPostTime time.Time `json:"last_post_time"`
	TotalReplies int       `json:"total_replies"`
	TotalViews   int       `json:"total_views"`
	Score        int       `json:"score"`
	Votes        Votes     `json:"votes"`
	WebPath      string    `json:"web_path"`
	Offer        *Offer    `json:"offer"`
	Image        string    `json:"image"`
}

type Votes struct {
	TotalUp   int `json:"total_up"`
	TotalDown int `json:"total_down"`
}

type Offer struct {
	CategoryID int    `json:"category_id"`
	DealerName string `json:"dealer_name"`
	URL        string `json:"url"`
	Price      string `json:"price"`
	OrigPrice  string `json:"original_price"`
	SaleInfo   string `json:"sale_info"`
	Savings    string `json:"savings"`
}

type PostsResponse struct {
	Pager Pager  `json:"pager"`
	Posts []Post `json:"posts"`
}

type Post struct {
	PostID     int       `json:"post_id"`
	TopicID    int       `json:"topic_id"`
	AuthorID   int       `json:"author_id"`
	Number     int       `json:"number"`
	Body       string    `json:"body"`
	PostTime   time.Time `json:"post_time"`
	Title      string    `json:"title"`
	Votes      PostVotes `json:"votes"`
	AuthorName string
}

type UserResponse struct {
	User User `json:"user"`
}

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type PostVotes struct {
	TotalUp   int `json:"total_up"`
	TotalDown int `json:"total_down"`
}

func (t *Topic) DealerName() string {
	if t.Offer != nil {
		return t.Offer.DealerName
	}
	return ""
}

func (t *Topic) DealURL() string {
	if t.Offer != nil && t.Offer.URL != "" {
		return t.Offer.URL
	}
	return "https://forums.redflagdeals.com" + t.WebPath
}

func (t *Topic) Price() string {
	if t.Offer == nil {
		return ""
	}
	return t.Offer.Price
}

func (t *Topic) Savings() string {
	if t.Offer == nil {
		return ""
	}
	return t.Offer.Savings
}

var categoryNames = map[int]string{
	9:  "Computers & Electronics",
	10: "Home & Garden",
	11: "Automotive",
	12: "Food & Drink",
	13: "Entertainment",
	14: "Fashion & Apparel",
	15: "Travel",
	16: "Finance",
	17: "Phones & Telecom",
}

func (t *Topic) CategoryID() int {
	if t.Offer == nil {
		return 0
	}
	return t.Offer.CategoryID
}

func (t *Topic) CategoryName() string {
	id := t.CategoryID()
	if name, ok := categoryNames[id]; ok {
		return name
	}
	return ""
}
