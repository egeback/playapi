package parsers

import (
	"github.com/egeback/playapi/internal/models"
)

const nrWorkers = 100

//ParserInterface struct
type ParserInterface interface {
	GetSeasons(show models.Show) []models.Season
	GetShows() []models.Show
	GetShowsWithSeasons() []models.Show
	Name() string
	PostCheckShows([]models.Show) []models.Show
}

//GetSeasonsConcurent get seasons with channels. nrWorkers number of workers runs work in parallel
func GetSeasonsConcurent(p ParserInterface, shows []models.Show) []models.Show {
	jobs := make(chan models.Show, len(shows))
	results := make(chan models.Show, len(shows))

	for w := 0; w < nrWorkers; w++ {
		go worker(p, jobs, results)
	}

	for _, show := range shows {
		jobs <- show
	}

	newShows := make([]models.Show, 0, len(shows))
	for a := 1; a <= len(shows); a++ {
		show := <-results
		newShows = append(newShows, show)
	}
	close(jobs)

	return newShows
}

func worker(p ParserInterface, jobs <-chan models.Show, results chan<- models.Show) {
	for j := range jobs {
		j.Seasons = p.GetSeasons(j)
		j.Prossesed = true
		results <- j
	}
}

var parsers []ParserInterface

//All returns all parsers defined
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
