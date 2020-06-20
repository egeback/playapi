package models

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

var shows = make([]Show, 0)

//ShowsAll return shows with query
func ShowsAll(q string) ([]Show, error) {
	if q == "" {
		return shows, nil
	}
	as := []Show{}
	for k, s := range shows {
		if q == s.Slug {
			as = append(as, shows[k])
		}
	}
	return as, nil
}
