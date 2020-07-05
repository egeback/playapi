package svtplay

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/egeback/playapi/internal/models"
	"github.com/egeback/playapi/internal/parsers"
	"github.com/egeback/playapi/internal/utils"
)

// Parser struct
type Parser struct {
	parsers.ParserInterface
}

//CreateParser returns new tv4play parser
func CreateParser() parsers.ParserInterface {
	return Parser{}
}

// Name return name of parser
func (p Parser) Name() string {
	return "svtPlay"
}

//GetShowsWithSeasons using SeasonConcurent
func (p Parser) GetShowsWithSeasons() []models.Show {
	shows := p.GetShows()
	return parsers.GetSeasonsConcurent(p, shows, parsers.WorkerGetSeasons)
}

//Build a url to interact with api.svt.se
func (p Parser) getURL(operation string, hashValue string, variables map[string]interface{}) string {
	extensions := fmt.Sprintf(`{"persistedQuery": {"version": 1, "sha256Hash": "%s"}}`, hashValue)
	extensions = utils.Quote(extensions)

	b, err := json.Marshal(variables)
	if err != nil {
		log.Println(err)
	}

	vars := utils.Quote(string(b))

	return fmt.Sprintf("https://api.svt.se/contento/graphql?"+
		"ua=svtplaywebb-play-render-prod-client&"+
		"operationName=%s&"+
		"variables=%s&"+
		"extensions=%s", operation, vars, extensions)
}

// GetSeasons for given show
func (p Parser) GetSeasons(show *models.Show) []models.Season {
	slug := show.Slug

	variables := map[string]interface{}{"titleSlugs": []string{*slug}}

	downloadURL := p.getURL("TitlePage", "4122efcb63970216e0cfb8abb25b74d1ba2bb7e780f438bbee19d92230d491c5", variables)
	result := utils.GetJSONFix(downloadURL)

	seasons := make([]models.Season, 0, 10)
	var data, ok = result["data"].(map[string]interface{})
	if !ok {
		log.Println("Could not convert result[\"data\"]")
		return seasons
	}
	var listablesBySlugContainer, ok2 = data["listablesBySlug"].([]interface{})
	if !ok2 {
		log.Println("Could not convert result[\"data\"]")
		return seasons
	}

	layout := "2006-01-02T15:04:05-07:00"
	latest, _ := time.Parse(layout, "1900-01-01T00:00:00+02:00")

	for _, l := range listablesBySlugContainer {
		listablesBySlugContainer := l.(map[string]interface{})
		associatedContent := listablesBySlugContainer["associatedContent"].([]interface{})
		for _, a := range associatedContent {
			content := a.(map[string]interface{})
			typeString := content["type"].(string)
			switch typeString {
			case "Season":
				episodes := make([]models.Episode, 0, 12)

				items := content["items"].([]interface{})
				for _, i := range items {
					itemContainer := i.(map[string]interface{})
					item := itemContainer["item"].(map[string]interface{})

					urls := item["urls"].(map[string]interface{})
					url := ""
					if x, ok := urls["svtplay"]; ok {
						url = "https://www.svtplay.se" + x.(string)
					}

					imageData := item["image"].(map[string]interface{})
					image := ""
					if x, ok := imageData["id"]; ok {
						image = fmt.Sprintf("https://www.svtstatic.se/image/large/1024/%s", x)
					}

					variants := make([]models.Variant, 0, 1)
					if x, ok := item["variants"]; ok {
						for _, v := range x.([]interface{}) {
							variantContainer := v.(map[string]interface{})
							url2 := ""
							if x, ok := urls["svtplay"]; ok {
								url2 = "https://www.svtplay.se" + x.(string)
							}
							platformSpecific := models.PlatformSpecific{"videoSvtId": variantContainer["videoSvtId"].(string)}
							variants = append(variants, models.Variant{ID: utils.GetStringValue(variantContainer, "id", nil), PlatformSpecific: platformSpecific, URL: &url2})
						}
					}

					id := utils.GetStringValue(item, "id", nil)
					name := utils.GetStringValue(item, "name", nil)
					svtID := utils.GetStringValue(item, "svtId", nil)
					videoSvtID := utils.GetStringValue(item, "videoSvtId", nil)
					slug := utils.GetStringValue(item, "slug", nil)
					descrption := utils.GetStringValue(item, "longDescription", nil)
					dur, ok := item["duration"]
					var duration float64 = 0
					if ok {
						duration = dur.(float64)
					} else {
					}
					num, ok := item["positionInSeason"]
					number := ""
					if ok {
						number = num.(string)
					} else {
					}
					vf, ok := item["validFrom"]
					validFromString := ""
					if ok {
						validFromString = vf.(string)
					}
					vt, ok := item["validTo"]
					validTo := ""
					if ok {
						validTo = vt.(string)
					}

					validFrom := utils.GetTimeFromString(validFromString)
					if validFrom != nil && latest.Before(*validFrom) {
						latest = *validFrom
					} else {
						if validFrom == nil {
							fmt.Println(show.Slug, validFromString)
						}
					}

					platformSpecific := models.PlatformSpecific{"svt_id": *svtID, "video_svt_id": *videoSvtID}
					episodes = append(episodes,
						models.Episode{
							ID: id, Name: name, Slug: slug, LongDescription: descrption,
							ImageURL: &image, URL: &url, Duration: &duration, Number: &number, ShowSlug: show.Slug, UpdatedAt: validFrom,
							ValidFrom: validFrom, ValidTo: utils.GetTimeFromString(validTo), Variants: variants, PlatformSpecific: platformSpecific})
				}

				season := models.Season{ID: utils.GetStringValue(content, "id", nil), Name: utils.GetStringValue(content, "name", nil), Episodes: episodes}
				seasons = append(seasons, season)
			case "Upcoming":
			case "Default":

			default:
				fmt.Println("No seasons for:", show.Name, show.Slug, typeString)
			}
		}
	}
	if latest.Year() != 1900 {
		show.UpdatedAt = &latest
	}

	return seasons
}

// GetShows from api.svt.se
func (p Parser) GetShows() []models.Show {
	gen := map[string]interface{}{"genre": []string{"dokumentar", "sport", "serier", "filmer", "barn", "drama", "humor",
		"reality", "underhållning", "samhalle-och-fakta", "kultur", "politik", "resor", "livsstil", "inspiration"}}
	//, "nyheter"
	url := p.getURL("GenreProgramsAO", "189b3613ec93e869feace9a379cca47d8b68b97b3f53c04163769dcffa509318", gen)
	result := utils.GetJSONFix(url)

	return p.extractShows(result)
}

func (p Parser) extractShows(result map[string]interface{}) []models.Show {
	shows := make([]models.Show, 0, 4000)
	var data = result["data"].(map[string]interface{})
	var genres = data["genres"].([]interface{})

	counter := 0

	for _, g := range genres {
		var genre = g.(map[string]interface{})
		var selectionsForWeb = genre["selectionsForWeb"].([]interface{})
		for z := 0; z < len(selectionsForWeb); z++ {
			selectionForWeb := selectionsForWeb[z].(map[string]interface{})
			var items = selectionForWeb["items"].([]interface{})
			for _, i := range items {
				item := (i.(map[string]interface{}))["item"].(map[string]interface{})
				if item["__typename"] == "Episode" || item["__typename"] == "Clip" || item["__typename"] == "Trailer" {
					if parent := utils.GetMapValue(item, "parent"); parent != nil && len(*parent) > 0 {
						if *utils.GetStringValue(item, "__typename", nil) != "Episode" {
							fmt.Println(*utils.GetStringValue(item, "name", nil), *utils.GetStringValue(item, "slug", nil))
							fmt.Println("	", *utils.GetStringValue(*parent, "name", nil))
							fmt.Println(*utils.GetStringValue(item, "__typename", nil))
						}
					} else if item["__typename"] != "Clip" && item["__typename"] != "Trailer" {
						fmt.Println("No parent:", *utils.GetStringValue(item, "name", nil), *utils.GetStringValue(item, "slug", nil))
					}
					continue
				} else if item["__typename"] == "Single" {
				} else if utils.Contains([]string{"KidsTvShow", "TvShow", "TvSeries"}, item["__typename"].(string)) {
					urls := item["urls"].(map[string]interface{})
					url := ""
					if x, ok := urls["svtplay"]; ok {
						url = "https://www.svtplay.se" + x.(string)
					}

					imageData := item["image"].(map[string]interface{})
					image := ""
					if x, ok := imageData["id"]; ok {
						image = fmt.Sprintf("https://www.svtstatic.se/image/large/1024/%s", x)
					}

					variables := map[string]interface{}{"titleSlugs": []string{item["slug"].(string)}}
					downloadURL := p.getURL("TitlePage", "4122efcb63970216e0cfb8abb25b74d1ba2bb7e780f438bbee19d92230d491c5", variables)
					svtID := utils.GetStringValue(item, "svtId", nil)
					platformSpecific := models.PlatformSpecific{"svtId": svtID}
					show := models.Show{
						ID:               *utils.GetStringValue(item, "id", nil),
						Name:             utils.GetStringValue(item, "name", nil),
						Slug:             utils.GetStringValue(item, "slug", nil),
						APIURL:           &downloadURL,
						PageURL:          &url,
						ImageURL:         &image,
						Description:      utils.GetStringValue(item, "longDescription", nil),
						UpdatedAt:        nil,
						Genre:            utils.GetStringValue(genre, "name", nil),
						Prossesed:        false,
						Provider:         "svtplay",
						PlatformSpecific: &platformSpecific,
					}
					counter++
					shows = append(shows, show)
				} else {
					fmt.Println(item["__typename"])
				}
			}
		}
	}
	return shows
}

// PostCheckShows to remove the ones that should not be visible
func (p Parser) PostCheckShows(shows []models.Show) []models.Show {
	newShows := make([]models.Show, 0, len(shows))
	showsWithOutEpisodes := make([]models.Show, 0, 0)
	for _, show := range shows {
		if len(show.Seasons) == 0 {
			showsWithOutEpisodes = append(showsWithOutEpisodes, show)
		} else {
			newShows = append(newShows, show)
		}
	}

	newShows = append(newShows, parsers.GetSeasonsConcurent(p, showsWithOutEpisodes, workerGetExtra)...)

	return newShows
}

func workerGetExtra(p parsers.ParserInterface, jobs <-chan models.Show, results chan<- models.Show) {
	for j := range jobs {
		j.Seasons = getDataForSingleEpisode(&j)
		j.Prossesed = true
		results <- j
	}
}

func getDataForSingleEpisode(show *models.Show) []models.Season {
	var duration *float64 = nil
	var updatedAt *time.Time = nil
	var id *string = &show.ID
	platformSpecfic := *show.PlatformSpecific
	if platformSpecfic != nil {
		if svtID, ok := platformSpecfic["svtId"]; ok {
			url := "https://api.svt.se/video/" + string(*svtID.(*string))
			respons := utils.GetJSONFix(url)
			//Since this can only be used for non shows check response
			if msg, ok := respons["msg"]; ok {
				if strings.HasPrefix(msg.(string), "Only SvtId for an episode version") {
					return make([]models.Season, 0)
				}
			}

			duration = utils.GetFloat64Value(respons, "contentDuration", nil)
			if rights := utils.GetMapValue(respons, "rights"); rights != nil {
				updatedAtString := utils.GetStringValue(*rights, "validFrom", &parsers.EmptyString)
				updatedAt = utils.GetTimeFromString(*updatedAtString)
			}

		}
	}
	if updatedAt != nil {
		show.UpdatedAt = updatedAt
	} else {
		return make([]models.Season, 0)
	}
	episode := models.Episode{
		Name:             show.Name,
		Slug:             show.Slug,
		ShowSlug:         show.Slug,
		ImageURL:         show.ImageURL,
		URL:              show.PageURL,
		ID:               id,
		UpdatedAt:        updatedAt,
		PlatformSpecific: show.PlatformSpecific,
		Duration:         duration,
	}
	name := "Only one episode"
	seasons := []models.Season{{
		Name:     &name,
		ID:       &show.ID,
		Episodes: []models.Episode{episode},
	}}
	return seasons
}
