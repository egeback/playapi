package models

// Season handles season info
type Season struct {
	ID       string    `json:"id" groups:"api,seasons"`
	Name     string    `json:"name" groups:"api,seasons"`
	Episodes []Episode `json:"episodes" groups:"api,seasons"`
	//Data interface
}
