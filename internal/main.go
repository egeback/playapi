package main

import (
	"fmt"
	"sort"

	"github.com/egeback/play_media_api/internal/controllers"
	_ "github.com/egeback/play_media_api/internal/docs"
	"github.com/egeback/play_media_api/internal/models"
	"github.com/egeback/play_media_api/internal/parsers"
	"github.com/egeback/play_media_api/internal/parsers/svtplay"
	"github.com/egeback/play_media_api/internal/parsers/tv4play"
	"github.com/egeback/play_media_api/internal/version"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	jsonp "github.com/tomwei7/gin-jsonp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Play Media API
// @version 1.0
// @description API including SVT and TV4 Play

// @contact.name API Support
// @contact.url http://xxxx.xxx.xx
// @contact.email support@egeback.se

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api/v1/
func main() {
	fmt.Printf("Running Play Media API version: %s (%s)\n", version.BuildVersion, version.BuildTime)
	parsers.Set([]parsers.ParserInterface{new(svtplay.Parser), new(tv4play.Parser)})
	//parsers.Set([]parsers.ParserInterface{new(tv4play.Parser)})

	//r := gin.Default()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(jsonp.JsonP())
	c := controllers.NewController()
	v1 := r.Group("/api/v1")
	{
		shows := v1.Group("/shows")
		{
			shows.GET("/", c.ListShows)
			shows.GET("/:slug", c.ShowShow)
			shows.GET("/:slug/seasons", c.ShowShow)
		}
		common := v1.Group("/")
		{
			common.GET("ping", c.Ping)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go func() {
		gocron.Every(15).Minutes().Do(updateShows)
		<-gocron.Start()
	}()
	updateShows()

	r.Run(":8080")
}

func updateShows() {
	shows := make([]models.Show, 0)
	for _, parser := range parsers.All("") {
		s := parser.GetShowsWithSeasons()
		shows = append(shows, parser.PostCheckShows(s)...)
	}

	showsWithSeasons := 0
	showsWithNoSeasons := 0

	for _, show := range shows {
		if show.Name == "" {
			fmt.Println(show)
		}
		if len(show.Seasons) == 0 {
			showsWithNoSeasons++
		} else {
			showsWithSeasons++
		}
	}
	sort.SliceStable(shows, func(i, j int) bool {
		return shows[i].Slug < shows[j].Slug
	})
	models.ShowsSet(shows)
	fmt.Println("Shows with Seasons", showsWithSeasons)
	fmt.Println("Shows with no Seasons", showsWithNoSeasons)

	fmt.Println("Done")
}
