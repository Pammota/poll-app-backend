package models

type Vote struct {
	PollID                string `json:"poll_id"`
	UserID                string `json:"user_id"`
	UserIP                string `json:"user_ip"`
	UserUniqueBrowserHash string `json:"user_unique_browser_hash"`
	Option                string `json:"option"`
}

type Options struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type Poll struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	AdminPassword string    `json:"admin_password"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	AllowMultiple bool      `json:"allow_multiple"`
	Options       []Options `json:"options"`
}
