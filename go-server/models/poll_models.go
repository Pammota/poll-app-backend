package models

type Vote struct {
	ID                    string   `json:"id"`
	PollID                string   `json:"poll_id"`
	UserIP                string   `json:"user_ip"`
	UserUniqueBrowserHash string   `json:"user_unique_browser_hash"`
	HasMultiple           bool     `json:"has_multiple"`
	Option                string   `json:"option"`
	Options               []string `json:"options"`
}

type Options struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type Poll struct {
	ID                      string    `json:"id"`
	Question                string    `json:"question"`
	Slug                    string    `json:"slug"`
	Author                  string    `json:"author"`
	AuthorUniqueBrowserHash string    `json:"author_unique_browser_hash"`
	StartDate               string    `json:"start_date"`
	EndDate                 string    `json:"end_date"`
	AllowMultiple           bool      `json:"allow_multiple"`
	Options                 []Options `json:"options"`
}
