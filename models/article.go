package models

// article data
type PostData struct {
	Date       string `json:"date"`
	Post       string `json:"post"`
	CommentURL string `json:"comment_url"`
}

type CommentData struct {
	Date    string `json:"date"`
	Comment string `json:"comments"`
}
