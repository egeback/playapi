package dplay

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/egeback/playapi/internal/models"
	"github.com/egeback/playapi/internal/parsers"
	"github.com/egeback/playapi/internal/utils"
	"github.com/google/uuid"
)

// Parser struct
type Parser struct {
	uuid     string
	uuidTime time.Time
	cookies  []*http.Cookie
	parsers.ParserInterface
}

//Included struct
type Included struct {
	Images   map[string]map[string]interface{}
	Routes   map[string]map[string]interface{}
	Genres   map[string]map[string]interface{}
	Seasons  map[string]map[string]interface{}
	Channels map[string]map[string]interface{}
}

var float0 = float64(0)

//CreateParser returns new Dplay parser
func CreateParser() parsers.ParserInterface {
	uuid, uuidTime, cookies := getKey()
	return Parser{uuid: uuid, uuidTime: uuidTime, cookies: cookies}
}

func getKey() (string, time.Time, []*http.Cookie) {
	id := uuid.New()
	uuid := id.String()
	uuid = strings.ReplaceAll(uuid, "-", "")
	url := fmt.Sprintf("https://disco-api.dplay.se/token?realm=dplayse&deviceId=%s&shortlived=true", uuid)
	uuidTime := time.Now()

	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:77.0) Gecko/20100101 Firefox/77.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "sv-SE,sv;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("TE", "Trailers")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req) //send request
	if err != nil {
		log.Println(err)
		log.Println("Could not get cookies")
	}

	cookies := resp.Cookies()
	// output := "# Netscape HTTP Cookie File\n"
	// for _, cookie := range cookies {
	// 	expires, err := time.Parse("Tue, 2 Jul 2006 15:04:05 GMT", cookie.RawExpires)
	// 	timestamp := "1909253725"
	// 	if err == nil {
	// 		timestamp = strconv.FormatInt(expires.Unix(), 10)
	// 	}
	// 	fmt.Println(cookie.Value)
	// 	//
	// 	output = output + fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\t%s", "disco-api.dplay.se", "FALSE", cookie.Path, "TRUE", timestamp, cookie.Name, cookie.Value)
	// 	//disco-api.dplay.se	FALSE	/	TRUE	1909251346	st	eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJVU0VSSUQ6ZHBsYXlzZTo3Mjg5N2Q2ZC0xYWU2LTQyY2ItYjdmYy1jMTVmNjVhNTIyMTkiLCJqdGkiOiJ0b2tlbi0yYmI4NzhlMi05MWIxLTRjOTMtYTBmOS0xNTM4MDhhOTViZGEiLCJhbm9ueW1vdXMiOmZhbHNlLCJpYXQiOjE1OTM4Njc0Nzl9.PM5_vVtpeagtf4iBuqwxRKPAzRTMyjq47mzjzG8CbNY
	// }

	// defer resp.Body.Close()

	// bytes, err := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()
	// fmt.Println(string(bytes))
	// fmt.Println(url)

	// ioutil.WriteFile("cookies.txt", []byte(output), 0644)

	return uuid, uuidTime, cookies
}

func (p Parser) getData(url string) map[string]interface{} {
	if p.uuidTime.Before(time.Now().AddDate(0, 0, -1)) {
		uuid, uuidTime, cookies := getKey()
		p.uuid = uuid
		p.uuidTime = uuidTime
		p.cookies = cookies
	}
	return utils.GetJSON(url)
}

// GetSeasons for given show
func (p Parser) GetSeasons(show *models.Show) []models.Season {
	page := 1
	parametersBase := "decorators=viewingHistory&include=images,primaryChannel,show,season&filter[videoType]=EPISODE,LIVE,FOLLOW_UP,STANDALONE&filter[show.id]=%s&page[size]=100&page[number]=%d&sort=-seasonNumber,-episodeNumber,videoType,earliestPlayableStart"
	parameters := utils.Quote(fmt.Sprintf(parametersBase, show.ID, page))
	url := fmt.Sprintf("https://disco-api.dplay.se/content/videos?%s", parameters)

	response := utils.GetJSON(url, p.cookies...)

	meta := utils.GetMapValue(response, "meta")
	pages := utils.GetIntValue(*meta, "totalPages", 1)

	included := getImagesRoutesGenres(response)
	seasons := extractSeasons(response, show, included)

	for page = 2; page < pages+1; page++ {
		parameters := utils.Quote(fmt.Sprintf(parametersBase, show.ID, page))
		url := fmt.Sprintf("https://disco-api.dplay.se/content/videos?%s", parameters)
		//fmt.Println(url)

		response = utils.GetJSON(url, p.cookies...)
		if response["included"] == nil {
			fmt.Println("Error", url)
			continue
		}
		included = getImagesRoutesGenres(response)
		seasons = append(seasons, extractSeasons(response, show, included)...)
	}

	return seasons
}

// GetShows from dplay
func (p Parser) GetShows() []models.Show {
	//url := fmt.Sprintf("https://%s/content/shows?include=images&page%5Bsize%5D=100&page%5Bnumber%5D={{0}}", baseURL)
	page := 1
	parametersBase := "include=images,genres,seasons&page[size]=100&page[number]=%d"
	parameters := utils.Quote(fmt.Sprintf(parametersBase, page))
	url := fmt.Sprintf("https://disco-api.dplay.se/content/shows?%s", parameters)

	response := utils.GetJSON(url, p.cookies...)

	meta := utils.GetMapValue(response, "meta")
	pages := utils.GetIntValue(*meta, "totalPages", 1)

	included := getImagesRoutesGenres(response)
	shows := extractShows(response, included)

	for page = 2; page < pages+1; page++ {
		parameters := utils.Quote(fmt.Sprintf(parametersBase, page))
		url := fmt.Sprintf("https://disco-api.dplay.se/content/shows?%s", parameters)

		response = utils.GetJSON(url, p.cookies...)
		included = getImagesRoutesGenres(response)
		shows = append(shows, extractShows(response, included)...)
	}

	return shows
}

func getImagesRoutesGenres(response map[string]interface{}) Included {
	incl := Included{
		Images:   make(map[string]map[string]interface{}),
		Genres:   make(map[string]map[string]interface{}),
		Routes:   make(map[string]map[string]interface{}),
		Seasons:  make(map[string]map[string]interface{}),
		Channels: make(map[string]map[string]interface{}),
	}

	for _, i := range response["included"].([]interface{}) {
		included := i.(map[string]interface{})

		t := *utils.GetStringValue(included, "type", &parsers.EmptyString)
		id := *utils.GetStringValue(included, "id", &parsers.EmptyString)
		attributes := utils.GetMapValue(included, "attributes")

		f := false
		if t == "image" {
			image := make(map[string]interface{})

			image["src"] = *utils.GetStringValue(*attributes, "src", &parsers.EmptyString)
			image["kind"] = *utils.GetStringValue(*attributes, "kind", &parsers.EmptyString)
			image["width"] = utils.GetIntValue(*attributes, "width", 0)
			image["height"] = utils.GetIntValue(*attributes, "height", 0)
			image["default"] = *utils.GetBoolValue(*attributes, "default", &f)
			image["id"] = id

			incl.Images[id] = image
		} else if t == "route" {
			route := make(map[string]interface{})

			route["url"] = *utils.GetStringValue(*attributes, "url", &parsers.EmptyString)
			route["canonical"] = *utils.GetBoolValue(*attributes, "canonical", &f)
			route["id"] = id
			incl.Routes[id] = route
		} else if t == "genre" {
			genre := make(map[string]interface{})

			genre["name"] = *utils.GetStringValue(*attributes, "name", &parsers.EmptyString)
			genre["alternateId"] = *utils.GetStringValue(*attributes, "alternateId", &parsers.EmptyString)
			genre["id"] = id
			incl.Genres[id] = genre
		} else if t == "channel" {
			channel := make(map[string]interface{})

			channel["name"] = *utils.GetStringValue(*attributes, "name", &parsers.EmptyString)
			channel["alternateId"] = *utils.GetStringValue(*attributes, "alternateId", &parsers.EmptyString)
			channel["description"] = *utils.GetStringValue(*attributes, "description", &parsers.EmptyString)
			channel["id"] = id
			incl.Channels[id] = channel
		} else if t == "season" {
			season := make(map[string]interface{})

			season["episodeCount"] = utils.GetIntValue(*attributes, "episodeCount", 0)
			season["plannedEpisodeCount"] = utils.GetIntValue(*attributes, "plannedEpisodeCount", 0)
			season["seasonNumber"] = utils.GetIntValue(*attributes, "seasonNumber", 0)
			season["videoCount"] = utils.GetIntValue(*attributes, "videoCount", 0)
			season["id"] = id
			incl.Seasons[id] = season
		}
	}
	return incl
}

func extractShows(response map[string]interface{}, included Included) []models.Show {
	data := response["data"].([]interface{})
	shows := make([]models.Show, 0, len(data))
	for _, d := range data {
		data := d.(map[string]interface{})

		attributes := utils.GetMapValue(data, "attributes")
		relationships := utils.GetMapValue(data, "relationships")

		id := utils.GetStringValue(data, "id", &parsers.EmptyString)
		description := utils.GetStringValue(*attributes, "description", &parsers.EmptyString)
		name := utils.GetStringValue(*attributes, "name", &parsers.EmptyString)
		newestEpisodePublishStartString := utils.GetStringValue(*attributes, "newestEpisodePublishStart", &parsers.EmptyString)
		newestEpisodePublishStart := utils.GetTimeFromString(*newestEpisodePublishStartString)
		slug := utils.GetStringValue(*attributes, "alternateId", &parsers.EmptyString)

		platformSpecific := models.PlatformSpecific{}

		imageURL := ""
		imageURLs := make([]interface{}, 0)
		if value, ok := (*utils.GetMapValue(*relationships, "images"))["data"]; ok {
			imgs := value.([]interface{})
			for _, i := range imgs {
				image := i.(map[string]interface{})
				i := included.Images[image["id"].(string)]

				if i["kind"].(string) == "logo" {
					imageURL = i["src"].(string)
				}
				imageURLs = append(imageURLs, i)
			}
			platformSpecific["images"] = imageURLs
		}

		routeData := make([]interface{}, 0)
		if value, ok := (*utils.GetMapValue(*relationships, "routes"))["data"]; ok {
			ruts := value.([]interface{})
			for _, i := range ruts {
				route := i.(map[string]interface{})
				routeData = append(routeData, included.Routes[route["id"].(string)])
			}
			platformSpecific["routes"] = routeData
		} else {
			fmt.Println("Empty routes", *name, *id)
		}

		gen := ""
		genreData := make([]interface{}, 0)
		if value, ok := (*utils.GetMapValue(*relationships, "genres"))["data"]; ok {
			gens := value.([]interface{})
			for _, i := range gens {
				genre := i.(map[string]interface{})
				genreData = append(genreData, included.Genres[genre["id"].(string)])
				gen = (included.Genres[genre["id"].(string)])["name"].(string)
			}
			platformSpecific["genres"] = genreData
		}

		seasonData := make([]interface{}, 0)
		if value, ok := (*utils.GetMapValue(*relationships, "routes"))["data"]; ok {
			seasons := value.([]interface{})
			for _, i := range seasons {
				season := i.(map[string]interface{})
				seasonData = append(routeData, included.Routes[season["id"].(string)])
			}
			platformSpecific["seasons"] = seasonData
		}

		pageURL := fmt.Sprintf("https://www.dplay.se/program/%s", *slug)
		apiURL := fmt.Sprintf("https://disco-api.dplay.se/content/videos?decorators=viewingHistory&include=images,primaryChannel,show&filter[videoType]=EPISODE,LIVE,FOLLOW_UP,STANDALONE&filter[show.id]=%s&page[size]=100&page[number]=1&sort=-seasonNumber,-episodeNumber,videoType,earliestPlayableStart", *id)

		show := models.Show{
			ID:               *id,
			Name:             name,
			Description:      description,
			UpdatedAt:        newestEpisodePublishStart,
			Slug:             slug,
			Genre:            &gen,
			PlatformSpecific: &platformSpecific,
			PageURL:          &pageURL,
			ImageURL:         &imageURL,
			Provider:         "dplay",
			APIURL:           &apiURL,
		}
		shows = append(shows, show)
	}

	return shows
}

func extractSeasons(response map[string]interface{}, show *models.Show, included Included) []models.Season {
	data := response["data"].([]interface{})
	seasons := make(map[int]models.Season)
	for _, d := range data {
		data := d.(map[string]interface{})

		attributes := utils.GetMapValue(data, "attributes")
		relationships := utils.GetMapValue(data, "relationships")
		id := utils.GetStringValue(data, "id", &parsers.EmptyString)

		publishStartString := utils.GetStringValue(*attributes, "publishStart", &parsers.EmptyString)
		publishStart := utils.GetTimeFromString(*publishStartString)

		validFrom := publishStart
		var validTo *time.Time = nil

		if availabilityWindows, exists := (*attributes)["availabilityWindows"]; exists {
			registered := false
			for _, aW := range availabilityWindows.([]interface{}) {
				availabilityWindow := aW.(map[string]interface{})
				validFrom = utils.GetTimeFromString(availabilityWindow["playableStart"].(string))
				validTo = utils.GetTimeFromString(availabilityWindow["playableStart"].(string))
				if availabilityWindow["package"] == "Registered" {
					registered = true
				}
			}
			if !registered {
				//continue
			}
		}

		slug := utils.GetStringValue(*attributes, "alternateId", &parsers.EmptyString)
		description := utils.GetStringValue(*attributes, "description", &parsers.EmptyString)
		name := utils.GetStringValue(*attributes, "name", &parsers.EmptyString)
		duration := utils.GetFloat64Value(*attributes, "videoDuration", &float0)
		path := utils.GetStringValue(*attributes, "path", &parsers.EmptyString)

		airDateString := utils.GetStringValue(*attributes, "airDate", &parsers.EmptyString)
		var airDate *time.Time
		if *airDateString != "" {
			airDate = utils.GetTimeFromString(*airDateString)
		} else {
			airDate = publishStart
		}

		seasonNumber := utils.GetIntValue(*attributes, "seasonNumber", 0)
		episodeNumber := utils.GetIntValue(*attributes, "episodeNumber", 0)

		var seasonData *map[string]interface{} = nil
		if value, ok := (*utils.GetMapValue(*relationships, "season"))["data"]; ok {
			s := value.(map[string]interface{})
			seasonData = &s
		}

		// if _, ok := included.Seasons[string(seasonNumber)]; !ok {
		// 	id := "0"
		// 	if seasonData["id"] != nil {
		// 		id = seasonData["id"].(string)
		// 	}
		// 	name := strconv.FormatInt(int64(seasonNumber), 10)
		// 	seasons[seasonNumber] = models.Season{
		// 		ID:       &id,
		// 		Name:     &name,
		// 		Episodes: make([]models.Episode, 0, 0),
		// 	}
		// }

		url := fmt.Sprintf("https://www.dplay.se/videos%s", *path)
		routeData := make(map[string]interface{})
		if value, ok := (*utils.GetMapValue(*relationships, "routes"))["data"]; ok {
			routes := value.([]interface{})
			for _, i := range routes {
				route := i.(map[string]interface{})
				routeData = route
			}
		}

		if len(routeData) > 0 {
			id := routeData["id"].(string)
			if route, ok := included.Routes[id]; ok {
				url = fmt.Sprintf("https://www.dplay.se/videos%s", route["url"].(string))
			}
		}

		imageURL := ""
		imageData := make(map[string]interface{})
		if value, ok := (*utils.GetMapValue(*relationships, "images"))["data"]; ok {
			images := value.([]interface{})
			for _, i := range images {
				image := i.(map[string]interface{})
				imageData = image
			}
		}

		if len(imageData) > 0 {
			id := imageData["id"].(string)
			if image, ok := included.Images[id]; ok {
				url = image["src"].(string)
			}
		}

		episodeNumberString := string(episodeNumber)
		episode := models.Episode{
			ID:              id,
			Name:            name,
			LongDescription: description,
			ValidFrom:       validFrom,
			ValidTo:         validTo,
			Duration:        duration,
			ImageURL:        &imageURL,
			Number:          &episodeNumberString,
			Slug:            slug,
			URL:             &url,
			UpdatedAt:       airDate,
		}
		if season, ok := seasons[seasonNumber]; ok {
			episodes := seasons[seasonNumber].Episodes
			season.Episodes = append(episodes, episode)
			seasons[seasonNumber] = season
		} else {
			id := "0"
			videoCount := 0
			if seasonData != nil && (*seasonData)["id"] != nil {
				id = (*seasonData)["id"].(string)

				if s, ok := included.Seasons[id]; ok {
					videoCount = s["videoCount"].(int)
				}
			}

			name := strconv.FormatInt(int64(seasonNumber), 10)
			episodes := make([]models.Episode, 0, videoCount)
			seasons[seasonNumber] = models.Season{
				ID:       &id,
				Name:     &name,
				Episodes: episodes,
			}
		}

	}
	values := make([]models.Season, 0, len(seasons))

	for _, value := range seasons {
		values = append(values, value)
	}

	return values
}

//GetShowsWithSeasons using SeasonConcurent
func (p Parser) GetShowsWithSeasons() []models.Show {
	shows := p.GetShows()
	return parsers.GetSeasonsConcurent(p, shows, parsers.WorkerGetSeasons)
}

//Name of service
func (p Parser) Name() string {
	return "dplay"
}

// PostCheckShows to remove the ones that should not be visible
func (p Parser) PostCheckShows(shows []models.Show) []models.Show {
	// for _, show := range shows {
	// 	if len(show.Seasons) > 0 {
	// 		fmt.Println(*show.Slug, *show.Name)
	// 	}
	// }
	return shows
}
