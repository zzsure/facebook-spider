package models

// article data
type ArticleData struct {
	Date     string   `json:"date"`
	Posts    string   `json:"content"`
	Comments []string `json:"comments"`
}
