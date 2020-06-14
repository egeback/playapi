package playmediaapi

// Show handles show information
type Show struct {
	ID          string   `json:"-" groups:"api"`
	Name        string   `json:"name" groups:"api"`
	Slug        string   `json:"slug" groups:"api"`
	URL         string   `json:"url" groups:"api"`
	Seasons     []Season `json:"seasons" groups:seasons"`
	ImageURL    string   `json:"imageUrl" groups:"api"`
	Description string   `json:"decription" groups:"api"`
	UpdatedAt   string   `json:"updatedAt" groups:"api"`
	Genre       string   `json:"genre" groups:"api"`
	//Data interface
}
