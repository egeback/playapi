package models

import (
	"strings"
	"time"

	"github.com/egeback/playapi/internal/utils"
)

// Show handles show information
type Show struct {
	ID               string            `json:"id" groups:"api" example:"1"`
	Name             *string           `json:"name" groups:"api" example:"Show Name"`
	Slug             *string           `json:"slug" groups:"api" example:"show_name"`
	APIURL           *string           `json:"api_url" groups:"api" example:"http://adad.ad/se"`
	PageURL          *string           `json:"page_url" groups:"api" example:"http://adad.ad/se"`
	Seasons          []Season          `json:"seasons" groups:"seasons"`
	ImageURL         *string           `json:"imageUrl" groups:"api" example:"http://adad.ad/se"`
	Description      *string           `json:"decription" groups:"api" example:"Show about x"`
	UpdatedAt        *time.Time        `json:"updatedAt" groups:"api" example:"2019-12-22"`
	Genre            *string           `json:"genre" groups:"api" example:"2019-12-22"`
	Prossesed        bool              `json:"-"`
	Provider         string            `json:"service" groups:"api"`
	PlatformSpecific *PlatformSpecific `json:"platform_specific" groups:"api"`
}

var shows = make([]Show, 0)

//QueryItem struct with field and value data
type QueryItem struct {
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
	Operator string      `json:"operator"`
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
				if exclude(q, *s.Slug, q.Operator) {
					match = false
					break
				}
			} else if strings.ToLower(q.Field) == "genre" {
				if exclude(q, *s.Genre, q.Operator) {
					match = false
					break
				}
			} else if strings.ToLower(q.Field) == "name" {
				if exclude(q, *s.Name, q.Operator) {
					match = false
					break
				}
			} else if strings.ToLower(q.Field) == "provider" || strings.ToLower(q.Field) == "service" {
				if exclude(q, s.Provider, q.Operator) {
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

//AddShow show to array
func AddShow(show Show) {
	shows = append(shows, show)
}

// Exclude Function to determine an item should be excluded based on query item and value
func exclude(q QueryItem, value string, operator string) bool {
	switch q.Value.(type) {
	case []interface{}:
		if !utils.Contains(utils.ExtractStringSlice(q.Value.([]interface{})), value) {
			return true
		}
	case interface{}:
		if q.Operator == "is" || q.Operator == "" {
			if q.Value.(string) != value {
				return true
			}
		} else if q.Operator == "in" {
			// This is case sensitive
			if !strings.Contains(value, q.Value.(string)) {
				return true
			}
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
		_, exists := g[*show.Genre]
		if !exists {
			s := make([]*Show, 0, 0)
			g[*show.Genre] = s
		}
		s := g[*show.Genre]
		s = append(s, &show)

		_, exists = v[show.Provider]
		if !exists {
			s := make([]*Show, 0, 0)
			v[show.Provider] = s
		}
		s = g[*show.Genre]
		s = append(s, &show)
	}
}
