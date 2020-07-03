package models

import (
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
