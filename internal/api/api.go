package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/egeback/play_media_api/internal/models"
	"github.com/egeback/play_media_api/internal/parsers"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-version"
	"github.com/jasonlvhit/gocron"
	"github.com/liip/sheriff"
)

// RestAPI ...
type RestAPI struct {
	router  *gin.Engine
	parsers []parsers.ParserInterface
	Shows   []models.Show
	v1      *gin.RouterGroup
}

// Init ...
func Init(parsers []parsers.ParserInterface) *RestAPI {
	api := new(RestAPI)
	api.parsers = parsers
	api.router = gin.Default()
	api.updateShows()
	sort.SliceStable(api.Shows, func(i, j int) bool {
		return api.Shows[i].Slug < api.Shows[j].Slug
	})
	api.addRoutes()
	go func() {
		gocron.Every(15).Minutes().Do(api.updateShows)
		<-gocron.Start()
	}()
	return api
}

func (api *RestAPI) updateShows() {
	shows := make([]models.Show, 0)
	for _, parser := range api.parsers {
		shows = append(shows, parser.GetShowsWithSeasons()...)
	}

	api.Shows = shows
	showsWithSeasons := 0
	showsWithNoSeasons := 0
	for _, show := range api.Shows {
		if show.Name == "" {
			fmt.Println(show)
		}
		if len(show.Seasons) == 0 {
			//fmt.Println(show.Name, show.Slug)
			showsWithNoSeasons++
		} else {
			showsWithSeasons++
		}
	}
	fmt.Println("Shows with Seasons", showsWithSeasons)
	fmt.Println("Shows with no Seasons", showsWithNoSeasons)

	fmt.Println("Done")
}

func (api *RestAPI) addRoutes() {
	api.v1 = api.router.Group("/api/v1")
	{
		api.v1.GET("/ping", api.ping)
		api.v1.GET("/shows", api.shows)
		api.v1.GET("/show/:slug", api.show)
		api.v1.GET("/shows/:slug", api.show)
	}
}

func (api RestAPI) marshalShows(version *version.Version, groups []string, prettyPrint bool) ([]byte, error) {
	groups = append(groups, "api")
	o := &sheriff.Options{
		Groups:     groups,
		ApiVersion: version,
	}

	data, err := sheriff.Marshal(o, api.Shows)
	if err != nil {
		return nil, err
	}
	if prettyPrint {
		return json.MarshalIndent(data, "", "  ")
	} else {
		return json.Marshal(data)
	}
}

// Start ...
func (api RestAPI) Start() {
	api.router.Run(":8080")
}

func (api RestAPI) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"data":       "pong",
		"status":     "ok",
		"data_type":  "String",
		"statusCode": 200,
	})
}
func (api RestAPI) shows(c *gin.Context) {
	seasons := c.DefaultQuery("seasons", "false")
	displaySeasons, err := strconv.ParseBool(seasons)
	pp := c.DefaultQuery("prettyPrint", "false")
	prettyPrint, err2 := strconv.ParseBool(pp)
	if err != nil || err2 != nil {
		log.Panic(err)
		log.Panic(err2)
		c.JSON(200, gin.H{
			"data":       "coudld not parse quyery parameters",
			"status":     "error",
			"statusCode": 400,
			"data_type":  "String",
		})
	}
	if !displaySeasons {
		v1, err := version.NewVersion("1.0.0")
		if err != nil {
			log.Panic(err)
		}
		output, err := api.marshalShows(v1, []string{"api"}, prettyPrint)

		c.Data(http.StatusOK, "application/json", output)
	} else {
		c.JSON(200, gin.H{
			"data":       api.Shows,
			"status":     "ok",
			"statusCode": 200,
			"data_type":  "Shows",
		})
	}
}

func (api RestAPI) show(c *gin.Context) {
	slug := c.Param("slug")
	prettyPrint := c.DefaultQuery("prettyPrint", "false")
	var show *models.Show

	for _, s := range api.Shows {
		if s.Slug == slug {
			show = &s
			break
		}
	}

	if show != nil {
		if strings.ToLower(prettyPrint) == "false" {
			c.JSON(200, gin.H{
				"data":       show,
				"status":     "ok",
				"statusCode": 200,
				"data_type":  "Show",
			})
			return
		}
		c.IndentedJSON(200, gin.H{
			"data":       show,
			"status":     "ok",
			"statusCode": 200,
			"data_type":  "Show",
		})
		return
	}
	c.JSON(200, gin.H{
		"data":       "Show not found",
		"status":     "not found",
		"statusCode": 404,
		"data_type":  "String",
	})
}
