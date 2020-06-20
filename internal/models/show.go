package models

// Show handles show information
type Show struct {
	ID          string   `json:"-" groups:"api" example:"1"`
	Name        string   `json:"name" groups:"api" example:"Show Name"`
	Slug        string   `json:"slug" groups:"api" example:"show_name"`
	URL         string   `json:"url" groups:"api" example:"http://adad.ad/se"`
	Seasons     []Season `json:"seasons" groups:"seasons"`
	ImageURL    string   `json:"imageUrl" groups:"api" example:"http://adad.ad/se"`
	Description string   `json:"decription" groups:"api" example:"Show about x"`
	UpdatedAt   string   `json:"updatedAt" groups:"api" example:"2019-12-22"`
	Genre       string   `json:"genre" groups:"api" example:"2019-12-22"`
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

//ShowsSet set shows
func ShowsSet(s []Show) {
	shows = s
}
