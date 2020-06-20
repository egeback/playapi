package parsers

import "github.com/egeback/play_media_api/internal/models"

const nrWorkers = 100

//ParserInterface ...
type ParserInterface interface {
	GetSeasons(show models.Show) []models.Season
	GetShows() []models.Show
	GetShowsWithSeasons() []models.Show
	Name() string
	//getURL(operation string, hashValue string, variables map[string]interface{}) string
	//GetSeasonsConcurent(shows []models.Show) []models.Show
}

//GetSeasonsConcurent ...
func GetSeasonsConcurent(p ParserInterface, shows []models.Show) []models.Show {
	jobs := make(chan models.Show, len(shows))
	results := make(chan models.Show, len(shows))

	for w := 0; w < nrWorkers; w++ {
		go worker(p, jobs, results)
	}

	for _, show := range shows {
		jobs <- show
	}

	close(jobs)
	newShows := make([]models.Show, 0, len(shows))
	for i := 0; i < len(shows); i++ {
		newShows = append(newShows, <-results)
	}
	return newShows
}

//Worker ...
func worker(p ParserInterface, jobs <-chan models.Show, results chan<- models.Show) {
	for j := range jobs {
		j.Seasons = p.GetSeasons(j)
		results <- j
	}
}

var parsers []ParserInterface

//All return all parsers defined
func All(q string) []ParserInterface {
	if q == "" {
		return parsers
	}
	as := []ParserInterface{}
	for k, p := range parsers {
		if q == p.Name() {
			as = append(as, parsers[k])
		}
	}
	return as
}

//Set set allication parsers
func Set(p []ParserInterface) {
	parsers = p
}
