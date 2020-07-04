package models

import (
	"strings"
	"time"

	"github.com/egeback/playapi/internal/utils"
)

// Episode is
type Episode struct {
	ID               *string     `json:"id" groups:"api"`
	Name             *string     `json:"name" groups:"api"`
	Slug             *string     `json:"slug" groups:"api"`
	LongDescription  *string     `json:"longDescription" groups:"api"`
	ImageURL         *string     `json:"imageUrl" groups:"api"`
	URL              *string     `json:"url" groups:"api"`
	Duration         *float64    `json:"duration" groups:"api"`
	Number           *string     `json:"number" groups:"api"`
	ValidFrom        *time.Time  `json:"validFrom" groups:"api"`
	ValidTo          *time.Time  `json:"validTo" groups:"api"`
	Variants         []Variant   `json:"variants" groups:"api"`
	PlatformSpecific interface{} `json:"platform_specific" groups:"api"`
	ShowSlug         *string     `json:"showSlug" groups:"api"`
	UpdatedAt        *time.Time  `json:"updated_at" groups:"api"`
}

//PlatformSpecific interface
type PlatformSpecific map[string]interface{}

var latestEpisodes = make(map[string][]Episode)
var allEpisodes = make([]Episode, 0)

//GetLatestEpisodes returns latest shows from the stream services
func GetLatestEpisodes(services ...string) ([]Episode, error) {
	if len(services) == 0 {
		episodes := make([]Episode, 0, 0)
		for _, e := range latestEpisodes {
			episodes = append(episodes, e...)
		}
		return episodes, nil
	}

	episodes := make([]Episode, 0, 0)
	for _, service := range services {
		if utils.Contains(services, service) {
			episodes = append(episodes, latestEpisodes[service]...)
		}
	}
	return episodes, nil
}

//ClearLatestEpisodes clear latest episodes
func ClearLatestEpisodes() {
	latestEpisodes = make(map[string][]Episode)
}

//AddLatestEpisodes add episodes
func AddLatestEpisodes(episodes []Episode, service string) {
	latestEpisodes[service] = episodes
}

//RemoveServiceFromLatestEpisodes remove from map
func RemoveServiceFromLatestEpisodes(service string) {
	delete(latestEpisodes, service)
}

//GetEpisodes returns latest shows from the stream services
func GetEpisodes(queryItems ...QueryItem) ([]Episode, error) {
	if len(queryItems) == 0 {
		return allEpisodes, nil
	}
	as := make([]Episode, 0, 0)
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
			as = append(as, allEpisodes[k])
		}
	}
	return as, nil
}

//ClearEpisodes clear latest episodes
func ClearEpisodes() {
	allEpisodes = make([]Episode, 0, 0)
}

//AddEpisodes add episodes
func AddEpisodes(episodes []Episode) {
	allEpisodes = append(allEpisodes, episodes...)
}
