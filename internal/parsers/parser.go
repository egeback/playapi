package parsers

import (
	"log"
	"sort"
	"time"

	"github.com/egeback/playapi/internal/models"
)

const nrWorkers = 100

//EmptyString to use as pointer
var EmptyString string = ""

//ParserInterface struct
type ParserInterface interface {
	GetSeasons(show *models.Show) []models.Season
	GetShows() []models.Show
	GetShowsWithSeasons() []models.Show
	//GetLatest(shows []models.Show) []models.Episode
	Name() string
	PostCheckShows([]models.Show) []models.Show
}

//GetSeasonsConcurent get seasons with channels. nrWorkers number of workers runs work in parallel
func GetSeasonsConcurent(p ParserInterface, shows []models.Show, fn worker) []models.Show {
	jobs := make(chan models.Show, len(shows))
	results := make(chan models.Show, len(shows))

	for w := 0; w < nrWorkers; w++ {
		go fn(p, jobs, results)
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

type worker func(ParserInterface, <-chan models.Show, chan<- models.Show)

//WorkerGetSeasons used to get Seasons on a given show
func WorkerGetSeasons(p ParserInterface, jobs <-chan models.Show, results chan<- models.Show) {
	for j := range jobs {
		j.Seasons = p.GetSeasons(&j)
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

//GetLatest from Service
func GetLatest(p ParserInterface, shows []models.Show) []models.Episode {
	if len(shows) == 0 {
		return make([]models.Episode, 0, 0)
	}

	episodes := make([]models.Episode, 0, 0)

	for _, show := range shows {
		if show.UpdatedAt == nil || (show.UpdatedAt != nil && time.Now().AddDate(0, 0, -7).Before(*show.UpdatedAt)) {
			episodes = checkEpisodes(show, episodes)
		}
	}

	//Sort slice
	sort.SliceStable(episodes, func(i, j int) bool {
		updatedAt := *episodes[i].UpdatedAt
		return updatedAt.Before(*episodes[j].UpdatedAt)
	})
	if len(episodes) <= 100 {
		return episodes
	}
	return episodes[:100]
}

func checkEpisodes(show models.Show, episodes []models.Episode) []models.Episode {
	for _, season := range show.Seasons {
		for _, episode := range season.Episodes {
			if episode.ValidFrom == nil {
				log.Println("No updated at for:", *show.Name, *show.Slug, len(show.Seasons))
			} else if time.Now().AddDate(0, 0, -7).Before(*episode.ValidFrom) {
				episodes = append(episodes, episode)
			}
		}
	}
	return episodes
}
