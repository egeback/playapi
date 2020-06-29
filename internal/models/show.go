package models

import (
	"strings"

	"github.com/egeback/playapi/internal/utils"
)

// Show handles show information
type Show struct {
	ID          string   `json:"-" groups:"api" example:"1"`
	Name        string   `json:"name" groups:"api" example:"Show Name"`
	Slug        string   `json:"slug" groups:"api" example:"show_name"`
	APIURL      string   `json:"api_url" groups:"api" example:"http://adad.ad/se"`
	PageURL     string   `json:"page_url" groups:"api" example:"http://adad.ad/se"`
	Seasons     []Season `json:"seasons" groups:"seasons"`
	ImageURL    string   `json:"imageUrl" groups:"api" example:"http://adad.ad/se"`
	Description string   `json:"decription" groups:"api" example:"Show about x"`
	UpdatedAt   string   `json:"updatedAt" groups:"api" example:"2019-12-22"`
	Genre       string   `json:"genre" groups:"api" example:"2019-12-22"`
	Prossesed   bool     `json:"-"`
	Provider    string   `json:"service" groups:"api"`
}

var shows = make([]Show, 0)

//QueryItem struct with field and value data
type QueryItem struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

//ShowsAll return shows with query, filters on slug, genre, name, provider/service
func ShowsAll(queryItems ...QueryItem) ([]Show, error) {
	if len(queryItems) == 0 {
		return shows, nil
	}
	as := []Show{}
	for k, s := range shows {
		match := true
		for _, q := range queryItems {
			if strings.ToLower(q.Field) == "slug" {
				if exclude(q, s.Slug) {
					match = false
					break
				}
			} else if strings.ToLower(q.Field) == "genre" {
				if exclude(q, s.Genre) {
					match = false
					break
				}
			} else if strings.ToLower(q.Field) == "name" {
				if exclude(q, s.Name) {
					match = false
					break
				}
			} else if strings.ToLower(q.Field) == "provider" || strings.ToLower(q.Field) == "service" {
				if exclude(q, s.Provider) {
					match = false
					break
				}
			}
		}
		if match {
			as = append(as, shows[k])
		}
	}
	return as, nil
}

// Function to determine an item should be excluded based on query item and value
func exclude(q QueryItem, value string) bool {
	switch q.Value.(type) {
	case []interface{}:
		if !utils.Contains(utils.ExtractStringSlice(q.Value.([]interface{})), value) {
			return true
		}
	case interface{}:
		if q.Value.(string) != value {
			return true
		}
	default:
		return true
	}
	return false
}

//ShowsSet set shows
func ShowsSet(s []Show) {
	shows = s
	g := make(map[string][]*Show, 0)
	v := make(map[string][]*Show, 0)
	for _, show := range shows {
		_, exists := g[show.Genre]
		if !exists {
			s := make([]*Show, 0, 0)
			g[show.Genre] = s
		}
		s := g[show.Genre]
		s = append(s, &show)

		_, exists = v[show.Provider]
		if !exists {
			s := make([]*Show, 0, 0)
			v[show.Provider] = s
		}
		s = g[show.Genre]
		s = append(s, &show)
	}
}
