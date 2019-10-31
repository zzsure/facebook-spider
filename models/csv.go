package models

// csv file data format
type FileData struct {
	URL  string `csv:"url"`      // official accounts url
	Lang string `csv:"language"` // official accounts language
}
