package parsers

import (
	"github.com/egeback/play_media_api/internal/models"
)

const nrWorkers = 100

//ParserInterface ...
type ParserInterface interface {
	GetSeasons(show models.Show) []models.Season
	GetShows() []models.Show
	GetShowsWithSeasons() []models.Show
	Name() string
	PostCheckShows([]models.Show) []models.Show
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

	newShows := make([]models.Show, 0, len(shows))
	for a := 1; a <= len(shows); a++ {
		//fmt.Printf("%d ", a)
		show := <-results
		newShows = append(newShows, show)
	}
	close(jobs)

	return newShows
	//
	// for i := 0; i < len(shows); i++ {
	// 	newShows = append(newShows, <-results)
	// }
	// return newShows
}

//GetSeasonsConcurent ...
// func GetSeasonsConcurent(p ParserInterface, shows []models.Show) []models.Show {
// 	l := len(shows)
// 	i := 0
// 	for ok := true; ok; ok = (i < l) {
// 		if i+nrWorkers < l {
// 			executeRange(p, shows[i:i+nrWorkers])
// 		} else {
// 			executeRange(p, shows[i:l])
// 			break
// 		}
// 		i += nrWorkers
// 	}
// 	//checkRange(shows)

// 	return shows
// }

// func checkRange(shows []models.Show) {
// 	count := 0
// 	s := make([]int, 0, 0)
// 	for i, show := range shows {
// 		if !show.Prossesed {
// 			count++
// 			s = append(s, i)
// 			fmt.Println(show.Name, show.Slug)
// 		}
// 	}
// 	fmt.Println("Nr of shows", len(shows), "not prossesed: ", count, s)
// }

// func executeRange(p ParserInterface, shows []models.Show) {
// 	results := make(chan bool, len(shows))
// 	for w := 0; w < len(shows); w++ {
// 		go worker(p, &shows[w], results)
// 	}
// 	for w := 0; w < len(shows); w++ {
// 		<-results
// 	}
// }

// func worker(p ParserInterface, show *models.Show, results chan bool) {
// 	show.Seasons = p.GetSeasons(*show)
// 	show.Prossesed = true
// 	results <- true

// }

//Worker ...
func worker(p ParserInterface, jobs <-chan models.Show, results chan<- models.Show) {
	for j := range jobs {
		j.Seasons = p.GetSeasons(j)
		j.Prossesed = true
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
