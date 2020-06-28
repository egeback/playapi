package tv4play

import (
	"fmt"

	"github.com/egeback/playapi/internal/models"
	"github.com/egeback/playapi/internal/parsers"
	"github.com/egeback/playapi/internal/utils"
	slugger "github.com/gosimple/slug"
)

// Parser ...
type Parser struct {
	parsers.ParserInterface
}

// Name return name of parser
func (p Parser) Name() string {
	return "tv4Play"
}

func updateGenre(genre string) string {
	updatedGenre := map[string]string{
		"dokumentärer":   "Dokumentär",
		"sport":          "Sport",
		"tv-serier":      "Serier",
		"nöje":           "Underhållning",
		"nyheter-debatt": "Samhälle & fakta",
	}
	value, exists := updatedGenre[genre]
	if exists {
		return value
	}
	return genre

}

func extractShow(data map[string]interface{}) models.Show {
	name := data["name"].(string)
	slug := data["nid"].(string)

	programImage := utils.GetStringValue(data, "program_image", "")
	id := utils.GetStringValue(data, "program_image", "")
	updatedAt := utils.GetStringValue(data, "updated_at", "")
	description := utils.GetStringValue(data, "description", "")
	genre := utils.GetStringValue(data, "category_nid", "")
	genre = updateGenre(genre)
	//videoAssetId := GetStringValue(data, "video_asset_id", "")

	url := fmt.Sprintf("https://api.tv4play.se/play/video_assets?platform=tablet&per_page=1000&is_live=false&type=episode&page=1&node_nids=%s&start=0", utils.Quote(slug))
	pageURL := fmt.Sprintf("https://www.tv4play.se/program/%s", utils.Quote(slug))

	return models.Show{
		ID:          id,
		Name:        name,
		Slug:        slug,
		APIURL:      url,
		PageURL:     pageURL,
		ImageURL:    programImage,
		Description: description,
		UpdatedAt:   updatedAt,
		Genre:       genre,
		Provider:    "tv4play",
	}
}

//GetSeasons ...
func (p Parser) GetSeasons(show models.Show) []models.Season {
	data := utils.GetJSON(show.APIURL)
	seasons := make(map[string]models.Season)

	var results = data["results"].([]interface{})

	if len(results) == 0 {
		return make([]models.Season, 0, 0)
	}

	for _, r := range results {
		result := r.(map[string]interface{})
		title := utils.GetStringValue(result, "title", "")
		slug := slugger.MakeLang(title, "sv")
		id := utils.GetStringValue(result, "id", "0")
		description := utils.GetStringValue(result, "description", "")
		//tags := GetStringValue(result, "tags", "")

		image := utils.GetStringValue(result, "image", "")
		season := utils.GetStringValue(result, "season", "0")
		episodeNr := utils.GetStringValue(result, "episode", "0")
		publishedDateTime := utils.GetStringValue(result, "published_date_time", "")
		expireDateTime := utils.GetStringValue(result, "expire_date_time", "")
		//isDrmProtected := GetStringValue(result, "is_drm_protected", "")

		_, exists := seasons[season]
		if !exists {
			s := models.Season{Name: season}
			s.Episodes = make([]models.Episode, 0, 1)
			seasons[season] = s
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
			URL:             fmt.Sprintf("https://www.tv4play.se/program/%s/%s", utils.Quote(slug), id),
		}
		s := seasons[season]
		s.Episodes = append(s.Episodes, episode)
		seasons[season] = s
	}

	values := make([]models.Season, 0, len(seasons))

	for _, value := range seasons {
		values = append(values, value)
	}

	return values
}

//GetShows ...
func (p Parser) GetShows() []models.Show {
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

// GetShowsWithSeasons ...
func (p Parser) GetShowsWithSeasons() []models.Show {
	shows := p.GetShows()
	return parsers.GetSeasonsConcurent(p, shows)
}

// PostCheckShows to remove the ones that should not be visible
func (p Parser) PostCheckShows(shows []models.Show) []models.Show {
	newShows := make([]models.Show, 0, len(shows))
	for _, show := range shows {
		if len(show.Seasons) > 0 {
			newShows = append(newShows, show)
		}
	}
	return newShows
}
