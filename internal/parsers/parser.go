package playmediaapi

import (
	playmediaapi "github.com/egeback/play_media_api"
)

const nrWorkers = 100

//ParserInterface ...
type ParserInterface interface {
	GetSeasons(show playmediaapi.Show) []playmediaapi.Season
	GetShows() []playmediaapi.Show
	GetShowsWithSeasons() []playmediaapi.Show
	//getURL(operation string, hashValue string, variables map[string]interface{}) string
	//GetSeasonsConcurent(shows []playmediaapi.Show) []playmediaapi.Show
}

//GetSeasonsConcurent ...
func GetSeasonsConcurent(p ParserInterface, shows []playmediaapi.Show) []playmediaapi.Show {
	jobs := make(chan playmediaapi.Show, len(shows))
	results := make(chan playmediaapi.Show, len(shows))

	for w := 0; w < nrWorkers; w++ {
		go worker(p, jobs, results)
	}

	for _, show := range shows {
		jobs <- show
	}

	close(jobs)
	newShows := make([]playmediaapi.Show, 0, len(shows))
	for i := 0; i < len(shows); i++ {
		newShows = append(newShows, <-results)
	}
	return newShows
}

//Worker ...
func worker(p ParserInterface, jobs <-chan playmediaapi.Show, results chan<- playmediaapi.Show) {
	for j := range jobs {
		j.Seasons = p.GetSeasons(j)
		results <- j
	}
}
