package tv4play

import (
	"fmt"
	"log"
	"time"

	"github.com/egeback/playapi/internal/models"
	"github.com/egeback/playapi/internal/parsers"
	"github.com/egeback/playapi/internal/utils"
	slugger "github.com/gosimple/slug"
)

// Parser stuct
type Parser struct {
	parsers.ParserInterface
	parsers.ParserBase
}

// Name return name of parser
func (p Parser) Name() string {
	return "tv4Play"
}

// Update genre to better match svtplay
func updateGenre(genre *string) *string {
	updatedGenre := map[string]string{
		"dokumentärer":   "Dokumentär",
		"sport":          "Sport",
		"tv-serier":      "Serier",
		"nöje":           "Underhållning",
		"nyheter-debatt": "Samhälle & fakta",
	}
	value, exists := updatedGenre[*genre]
	if exists {
		return &value
	}
	return genre

}

func extractShow(data map[string]interface{}) models.Show {
	name := utils.GetStringValue(data, "name", nil)
	slug := utils.GetStringValue(data, "nid", nil)

	programImage := utils.GetStringValue(data, "program_image", nil)
	updatedAtString := utils.GetStringValue(data, "updated_at", nil)
	description := utils.GetStringValue(data, "description", nil)
	genre := utils.GetStringValue(data, "category_nid", nil)
	genre = updateGenre(genre)

	url := fmt.Sprintf("https://api.tv4play.se/play/video_assets?platform=tablet&per_page=1000&is_live=false&type=episode&page=1&node_nids=%s&start=0", utils.Quote(*slug))
	pageURL := fmt.Sprintf("https://www.tv4play.se/program/%s", utils.Quote(*slug))

	updatedAt := utils.GetTimeFromString(*updatedAtString)

	return models.Show{
		//ID:          *id,
		Name:        name,
		Slug:        slug,
		APIURL:      &url,
		PageURL:     &pageURL,
		ImageURL:    programImage,
		Description: description,
		UpdatedAt:   updatedAt,
		Genre:       genre,
		Provider:    "tv4play",
	}
}

func (p Parser) getSeasonsBackup(show *models.Show) []models.Season {
	extensions := fmt.Sprintf(`{"persistedQuery": {"version": 1, "sha256Hash": "%s"}}`, "906eba9587ac10bd55ebb063b549b513bb690bc26dd373d95797efc57bedba67")
	extensions = utils.Quote(extensions)
	vars := utils.Quote(fmt.Sprintf("{\"nid\":\"%s\"}", *show.Slug))
	apiURL := fmt.Sprintf("https://graphql.tv4play.se/graphql?"+
		"operationName=%s&"+
		"variables=%s&"+
		"extensions=%s", "cdp", vars, extensions)

	response := utils.GetJSON(apiURL)
	data := utils.GetMapValue(response, "data")
	if (*data)["program"] == nil {
		return make([]models.Season, 0)
	}

	program := utils.GetMapValue(*data, "program")
	show.Description = utils.GetStringValue(*program, "decription", nil)

	seasons := make(map[string]models.Season)

	layout := "2006-01-02T15:04:05-07:00"
	layout2 := "2006-01-02T15:04:05Z"
	latest, _ := time.Parse(layout, "1900-01-01T00:00:00+02:00")

	pI := (*program)["panels"].([]interface{})

	for _, pO := range pI {
		panel := pO.(map[string]interface{})
		if *(utils.GetStringValue(panel, "assetType", &parsers.EmptyString)) != "EPISODE" {
			continue
		}

		videoList := utils.GetMapValue(panel, "videoList")
		for _, vA := range (*videoList)["videoAssets"].([]interface{}) {
			videoAsset := vA.(map[string]interface{})
			title := utils.GetStringValue(videoAsset, "title", nil)
			slug := slugger.MakeLang(*title, "sv")
			id := utils.GetStringValue(videoAsset, "id", nil)
			description := utils.GetStringValue(videoAsset, "description", nil)

			image := utils.GetStringValue(videoAsset, "image", nil)
			season := utils.GetStringValue(videoAsset, "season", nil)
			episodeNr := utils.GetStringValue(videoAsset, "episode", nil)
			publishedDateTimeString := utils.GetStringValue(videoAsset, "published_date_time", nil)
			broadcastDateTimeString := utils.GetStringValue(videoAsset, "broadcastDateTime", nil)
			expireDateTimeString := utils.GetStringValue(videoAsset, "expire_date_time", nil)

			updatedAtString := publishedDateTimeString
			var updatedAt *time.Time

			platformSpecific := models.PlatformSpecific{
				"broadcastDateTime": utils.GetTimeFromString(*broadcastDateTimeString),
				"seasonizedTitle":   utils.GetStringValue(videoAsset, "seasonizedTitle", nil),
				"geoRestricted":     utils.GetBoolValue(videoAsset, "geoRestricted", nil),
			}

			if updatedAtString != nil {
				updatedAtTime := utils.GetTimeFromString(*updatedAtString)
				if updatedAtTime != nil {
					updatedAt = updatedAtTime
					if latest.Before(*updatedAtTime) {
						latest = *updatedAtTime
					}
				}
			} else {
				fmt.Println("")
			}

			_, exists := seasons[*season]
			if !exists {
				s := models.Season{Name: season}
				s.Episodes = make([]models.Episode, 0, 1)
				seasons[*season] = s
			}
			url := fmt.Sprintf("https://www.tv4play.se/program/%s/%s", utils.Quote(slug), *id)
			var publishedDateTime *time.Time = nil
			if publishedDateTimeString != nil {
				publishedDateTime = utils.GetTimeFromString(*publishedDateTimeString, layout2)
			}

			var expireDateTime *time.Time = nil
			if expireDateTime != nil {
				expireDateTime = utils.GetTimeFromString(*expireDateTimeString, layout2)
			} else {
				daysLeftInService := utils.GetIntValue(videoAsset, "daysLeftInService", 0)
				expireDate := time.Now().AddDate(0, 0, daysLeftInService)
				expireDateTime = &expireDate
			}

			episode := models.Episode{
				ID:               id,
				Name:             title,
				Slug:             &slug,
				LongDescription:  description,
				ImageURL:         image,
				ValidFrom:        publishedDateTime,
				ValidTo:          expireDateTime,
				Number:           episodeNr,
				URL:              &url,
				ShowSlug:         show.Slug,
				UpdatedAt:        updatedAt,
				PlatformSpecific: platformSpecific,
				Duration:         utils.GetFloat64Value(videoAsset, "duration", nil),
			}
			s := seasons[*season]
			s.Episodes = append(s.Episodes, episode)
			seasons[*season] = s
		}
	}

	values := make([]models.Season, 0, len(seasons))

	for _, value := range seasons {
		values = append(values, value)
	}
	show.UpdatedAt = &latest
	return values
}

//GetSeasons from given show
func (p Parser) GetSeasons(show *models.Show) []models.Season {
	data := utils.GetJSON(*show.APIURL)
	seasons := make(map[string]models.Season)

	if _, ok := data["results"]; !ok {
		fmt.Print("No result for:", show.Name, show.Slug, show.APIURL)
	}

	var results = data["results"].([]interface{})

	if len(results) == 0 {
		return p.getSeasonsBackup(show)
	}

	layout := "2006-01-02T15:04:05-07:00"
	latest, _ := time.Parse(layout, "1900-01-01T00:00:00+02:00")

	for _, r := range results {
		result := r.(map[string]interface{})
		title := utils.GetStringValue(result, "title", nil)
		slug := utils.GetStringValue(result, "program_nid", nil)

		id := utils.GetStringValue(result, "id", nil)
		description := utils.GetStringValue(result, "description", nil)

		image := utils.GetStringValue(result, "image", nil)
		duration := utils.GetFloat64Value(result, "duration", nil)
		season := utils.GetStringValue(result, "season", nil)
		episodeNr := utils.GetStringValue(result, "episode", nil)
		publishedDateTimeString := utils.GetStringValue(result, "published_date_time", nil)
		expireDateTimeString := utils.GetStringValue(result, "expire_date_time", nil)

		updatedAtString := ""
		var updatedAt *time.Time
		if _, ok := result["program"]; ok {
			if result["program"] != nil {
				if program := utils.GetMapValue(result, "program"); program != nil {
					updatedAtString = *utils.GetStringValue(*program, "updated_at", nil)
					if show.Description == nil {
						show.Description = utils.GetStringValue(*program, "description", nil)
					}
				}
			} else {
				updatedAtString = *utils.GetStringValue(result, "broadcast_date_time", nil)
			}
		}

		if updatedAtString != "" {
			updatedAtTime, err := time.Parse(layout, updatedAtString)
			if err != nil {
				log.Println(err)
			} else {
				updatedAt = &updatedAtTime
				if latest.Before(updatedAtTime) {
					latest = updatedAtTime
				}
			}
		}
		if season == nil {
			s := "Season 1"
			season = &s
		}

		_, exists := seasons[*season]
		if !exists {
			s := models.Season{Name: season}
			s.Episodes = make([]models.Episode, 0, 1)
			seasons[*season] = s
		}
		url := ""
		if id != nil {
			url = fmt.Sprintf("https://www.tv4play.se/program/%s/%s", *show.Slug, *id)
		} else {
			url = fmt.Sprintf("https://www.tv4play.se/program/%s", *show.Slug)
		}
		//hd, is_drm_protected

		var publishedDateTime *time.Time = nil
		if publishedDateTimeString != nil {
			publishedDateTime = utils.GetTimeFromString(*publishedDateTimeString, layout)
		} else {
			log.Println(*show.Slug, *show.APIURL)
		}

		var expireDateTime *time.Time = nil
		if expireDateTimeString != nil {
			expireDateTime = utils.GetTimeFromString(*expireDateTimeString, layout)
		}

		episode := models.Episode{
			ID:              id,
			Name:            title,
			Slug:            slug,
			LongDescription: description,
			ImageURL:        image,
			ValidFrom:       publishedDateTime,
			ValidTo:         expireDateTime,
			Number:          episodeNr,
			URL:             &url,
			ShowSlug:        show.Slug,
			UpdatedAt:       updatedAt,
			Duration:        duration,
		}
		s := seasons[*season]
		s.Episodes = append(s.Episodes, episode)
		seasons[*season] = s
	}

	values := make([]models.Season, 0, len(seasons))

	for _, value := range seasons {
		values = append(values, value)
	}
	show.UpdatedAt = &latest
	return values
}

//GetShows from api.tv4play.se
func (p Parser) GetShows() []models.Show {
	response := utils.GetJSON("https://graphql.tv4play.se/graphql?operationName=ProgramSearch&variables=%7B%22order_by%22%3A%22NAME%22%2C%22per_page%22%3A1000%7D&extensions=%7B%22persistedQuery%22%3A%7B%22version%22%3A1%2C%22sha256Hash%22%3A%2278cdda0280f7e6b21dea52021406cc44ef0ce37102cb13571804b1a5bd3b9aa1%22%7D%7D")
	var data = response["data"].(map[string]interface{})
	var programSearch = data["programSearch"].(map[string]interface{})

	shows := make([]models.Show, 0, 0)

	for _, p := range programSearch["programs"].([]interface{}) {
		program := p.(map[string]interface{})
		name := utils.GetStringValue(program, "name", nil)
		slug := utils.GetStringValue(program, "nid", nil)

		programImage := utils.GetStringValue(program, "image", nil)
		id := utils.GetStringValue(program, "id", nil)
		category := utils.GetMapValue(program, "category")
		genre := utils.GetStringValue(*category, "name", nil)

		genre = updateGenre(genre)

		url := fmt.Sprintf("https://api.tv4play.se/play/video_assets?platform=tablet&per_page=1000&is_live=false&type=episode&page=1&node_nids=%s&start=0", utils.Quote(*slug))
		pageURL := fmt.Sprintf("https://www.tv4play.se/program/%s", utils.Quote(*slug))

		show := models.Show{
			ID:       *id,
			Name:     name,
			Slug:     slug,
			APIURL:   &url,
			PageURL:  &pageURL,
			ImageURL: programImage,
			Genre:    genre,
			Provider: "tv4play",
		}
		shows = append(shows, show)
	}
	return shows
}

//GetShows2 from api.tv4play.se
func (p Parser) GetShows2() []models.Show {
	data := utils.GetJSON("https://api.tv4play.se/play/programs?is_active=true&platform=tablet&per_page=1000&fl=nid,name,program_image,is_premium,updated_at,channel,description,category_nid&start=0")
	totalHits := int(data["total_hits"].(float64))
	var results = data["results"].([]interface{})

	shows := make([]models.Show, 0, totalHits)

	for _, r := range results {
		result := r.(map[string]interface{})
		shows = append(shows, extractShow(result))
	}

	return shows

}

// GetShowsWithSeasons by using GetSeasonsConcurent
func (p Parser) GetShowsWithSeasons() []models.Show {
	shows := p.GetShows()
	shows2 := p.GetShows2()

	showMap := make(map[string]models.Show)

	for _, show := range shows2 {
		if _, ok := showMap[*show.Slug]; !ok {
			showMap[*show.Slug] = show
		}
	}

	for _, show := range shows {
		if _, ok := showMap[*show.Slug]; !ok {
			showMap[*show.Slug] = show
		}
	}

	unionShows := make([]models.Show, 0, len(showMap))
	for _, show := range showMap {
		unionShows = append(unionShows, show)
	}
	//fmt.Println("Union of shows len", len(unionShows))

	return parsers.GetSeasonsConcurent(p, unionShows, parsers.WorkerGetSeasons)
}

// PostCheckShows to remove the ones that should not be visible
func (p Parser) PostCheckShows(shows []models.Show) []models.Show {
	newShows := make([]models.Show, 0, len(shows))
	for _, show := range shows {
		if len(show.Seasons) > 0 {
			newShows = append(newShows, show)
		} else {
			//fmt.Println(show.Name, show.Slug, show.APIURL)
		}
	}
	return newShows
}
