package svtplay

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/egeback/play_media_api/internal/models"
	"github.com/egeback/play_media_api/internal/parsers"
	"github.com/egeback/play_media_api/internal/utils"
)

// Parser ...
type Parser struct {
	parsers.ParserInterface
}

// Mapable is
type Mapable interface {
}

// Name return name of parser
func (p Parser) Name() string {
	return "svtPlay"
}

//GetShowsWithSeasons ...
func (p Parser) GetShowsWithSeasons() []models.Show {
	shows := p.GetShows()
	return parsers.GetSeasonsConcurent(p, shows)
}

func (p Parser) getURL(operation string, hashValue string, variables map[string]interface{}) string {
	extensions := fmt.Sprintf(`{"persistedQuery": {"version": 1, "sha256Hash": "%s"}}`, hashValue)
	extensions = utils.Quote(extensions)

	b, err := json.Marshal(variables)
	if err != nil {
		panic(err)
	}

	vars := utils.Quote(string(b))

	return fmt.Sprintf("https://api.svt.se/contento/graphql?"+
		"ua=svtplaywebb-play-render-prod-client&"+
		"operationName=%s&"+
		"variables=%s&"+
		"extensions=%s", operation, vars, extensions)
}

// GetSeasons ...
func (p Parser) GetSeasons(show models.Show) []models.Season {
	slug := show.Slug
	variables := map[string]interface{}{"titleSlugs": []string{slug}}

	downloadURL := p.getURL("TitlePage", "4122efcb63970216e0cfb8abb25b74d1ba2bb7e780f438bbee19d92230d491c5", variables)
	result := utils.GetJSON(downloadURL)

	//seasons := list.New()
	seasons := make([]models.Season, 0, 10)
	var data, ok = result["data"].(map[string]interface{})
	if !ok {
		log.Panic("Could not convert result[\"data\"]")
		return seasons
	}
	var listablesBySlugContainer, ok2 = data["listablesBySlug"].([]interface{})
	if !ok2 {
		log.Panic("Could not convert result[\"data\"]")
		return seasons
	}
	//var data = utils.GetInterfaceMap(result["data"])
	//var listablesBySlugContainer = utils.GetInterfaceArray(data["listablesBySlug"])

	for _, l := range listablesBySlugContainer {
		listablesBySlugContainer := l.(map[string]interface{})
		associatedContent := listablesBySlugContainer["associatedContent"].([]interface{})
		for _, a := range associatedContent {
			content := a.(map[string]interface{})
			typeString := content["type"].(string)
			switch typeString {
			case "Season":
				//episodes := list.New()
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

					//variants := list.New()
					variants := make([]models.Variant, 0, 1)
					if x, ok := item["variants"]; ok {
						for _, v := range x.([]interface{}) {
							variantContainer := v.(map[string]interface{})
							url2 := ""
							if x, ok := urls["svtplay"]; ok {
								url2 = "https://www.svtplay.se" + x.(string)
							}
							//variants.PushBack(
							platformSpecific := models.PlatformSpecific{"svt_id": variantContainer["videoSvtId"].(string)}
							variants = append(variants, models.Variant{ID: variantContainer["id"].(string), PlatformSpecific: platformSpecific, URL: url2})
						}
					}

					//episodes.PushBack(
					id := item["id"].(string)
					name := item["name"].(string)
					svtID := item["svtId"].(string)
					videoSvtID := item["videoSvtId"].(string)
					slug := item["slug"].(string)
					descrption := item["longDescription"].(string)
					dur, ok := item["duration"]
					var duration float64 = 0
					if ok {
						duration = dur.(float64)
					} else {
						//fmt.Println("Error extracting duration from:", show.Name, name, downloadURL)
					}
					num, ok := item["positionInSeason"]
					number := ""
					if ok {
						number = num.(string)
					} else {
						//fmt.Println("Error extracting duration from:", slug, name)
					}
					vf, ok := item["validFrom"]
					validFrom := ""
					if ok {
						validFrom = vf.(string)
					}
					vt, ok := item["validTo"]
					validTo := ""
					if ok {
						validTo = vt.(string)
					}
					platformSpecific := models.PlatformSpecific{"svt_id": svtID, "video_svt_id": videoSvtID}
					episodes = append(episodes,
						models.Episode{
							ID: id, Name: name, Slug: slug, LongDescription: descrption,
							ImageURL: image, URL: url, Duration: duration, Number: number,
							ValidFrom: validFrom, ValidTo: validTo, Variants: variants, PlatformSpecific: platformSpecific})
				}
				season := models.Season{content["id"].(string), content["name"].(string), episodes}
				//seasons.PushBack(season)
				seasons = append(seasons, season)
			case "Upcoming":
			case "Default":
				/*if len(seasons) == 0 {
					fmt.Println("No seasons for:", show.Name, show.Slug, typeString)
				}*/
			default:
				fmt.Println("No seasons for:", show.Name, show.Slug, typeString)
			}
		}
	}
	return seasons
}

// GetShows ...
func (p Parser) GetShows() []models.Show {
	gen := map[string]interface{}{"genre": []string{"dokumentar", "sport", "nyheter", "serier", "filmer", "barn", "drama", "humor",
		"reality", "underhållning", "samhalle-och-fakta", "kultur", "politik"}}

	url := p.getURL("GenreProgramsAO", "189b3613ec93e869feace9a379cca47d8b68b97b3f53c04163769dcffa509318", gen)
	result := utils.GetJSON(url)

	//shows := list.New()
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

				if item["slug"].(string) == "1-000-doda-pa-ett-dygn-i-brasilien" {
					//print(item)
					//fmt.Println(item["slug"].(string), url)
				}

				variables := map[string]interface{}{"titleSlugs": []string{item["slug"].(string)}}
				downloadURL := p.getURL("TitlePage", "4122efcb63970216e0cfb8abb25b74d1ba2bb7e780f438bbee19d92230d491c5", variables)

				show := models.Show{
					ID:          item["id"].(string),
					Name:        item["name"].(string),
					Slug:        item["slug"].(string),
					APIURL:      downloadURL,
					PageURL:     url,
					ImageURL:    image,
					Description: item["longDescription"].(string),
					UpdatedAt:   "",
					Genre:       genre["name"].(string),
					Prossesed:   false,
					Provider:    "svtplay",
				}
				//shows.PushBack(show)
				counter++
				shows = append(shows, show)
				/*if counter == 100 {
					return shows
				}*/
			}
		}
	}
	return shows
}

// PostCheckShows to remove the ones that should not be visible
func (p Parser) PostCheckShows(shows []models.Show) []models.Show {
	return shows
}
