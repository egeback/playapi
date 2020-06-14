package playmediaapi

// Season handles season info
type Season struct {
	ID       string    `json:"id groups:"api"`
	Name     string    `json:"name groups:"api"`
	Episodes []Episode `json:"episodes groups:"api"`
	//Data interface
}
