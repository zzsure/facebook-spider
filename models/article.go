package models

// article data
type ArticleData struct {
	Posts    string   `json:"content"`
	Comments []string `json:"comments"`
}
