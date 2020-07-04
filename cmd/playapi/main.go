package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/egeback/playapi/internal/controllers"
	_ "github.com/egeback/playapi/internal/docs"
	"github.com/egeback/playapi/internal/models"
	"github.com/egeback/playapi/internal/parsers"
	"github.com/egeback/playapi/internal/parsers/dplay"
	"github.com/egeback/playapi/internal/version"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	jsonp "github.com/tomwei7/gin-jsonp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/natefinch/lumberjack.v2"
)

// @title Play service API
// @version 1.0.5
// @description API including SVT and TV4 Play

// @contact.name API Support
// @contact.url https://github.com/egeback/playapi/issues
// @contact.email support@egeback.se

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api/v1/
func main() {
	// Configure Logging
	logFileLocation := os.Getenv("DOWNLOADER_LOG_FILE_LOCATION")
	if logFileLocation != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFileLocation,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Running Play Media API version: %s (%s)\n", version.BuildVersion, version.BuildTime)

	//Add parsers
	//parsers.Set([]parsers.ParserInterface{new(svtplay.Parser), new(tv4play.Parser)})
	parsers.Set([]parsers.ParserInterface{dplay.CreateParser()})

	//Configure gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(jsonp.JsonP())
	c := controllers.NewController()
	v1 := r.Group("/api/v1")
	{
		shows := v1.Group("/shows")
		{
			shows.GET("", c.ListShows)
			shows.GET("/", c.ListShows)
			shows.GET("/:slug", c.ShowShow)
			shows.GET("/:slug/seasons", c.ShowShow)
		}
		episodes := v1.Group("/episodes")
		{
			episodes.GET("", c.ListEpisodes)
			episodes.GET("/latest", c.ListLatestEpisodes)
		}
		common := v1.Group("/")
		{
			common.GET("/ping", c.Ping)
			common.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	//Schedule updates of shows
	go func() {
		gocron.Every(15).Minutes().Do(updateShows)
		<-gocron.Start()
	}()

	//Run inital show parsing
	updateShows()

	//Start server
	r.Run(":8080")
	os.Exit(0)
}

//Update shows by iterating parsers and run GetShowsWithSeasons
func updateShows() {
	t1 := time.Now()
	shows := make([]models.Show, 0)
	for _, parser := range parsers.All("") {
		s := parser.GetShowsWithSeasons()
		checkedShows := parser.PostCheckShows(s)
		shows = append(shows, checkedShows...)
		models.AddLatestEpisodes(parsers.GetLatest(parser, checkedShows), parser.Name())
	}

	showsWithSeasons := 0
	showsWithNoSeasons := 0
	notProcessed := 0

	//Gather statistics and add episodes
	episodes := make([]models.Episode, 0, 0)
	for _, show := range shows {
		for _, season := range show.Seasons {
			episodes = append(episodes, season.Episodes...)
		}
		if show.Name == nil {
			fmt.Println(show)
		}
		if len(show.Seasons) == 0 {
			showsWithNoSeasons++
		} else {
			showsWithSeasons++
		}
		if !show.Prossesed {
			notProcessed++
		}
	}
	//Sort episodes slice
	sort.SliceStable(episodes, func(i, j int) bool {
		if episodes[i].Slug != nil && episodes[j].Slug != nil {
			return *episodes[i].Slug < *episodes[j].Slug
		}
		return *episodes[i].Name < *episodes[j].Name
	})

	models.AddEpisodes(episodes)

	//Sort shows slice
	sort.SliceStable(shows, func(i, j int) bool {
		return *shows[i].Slug < *shows[j].Slug
	})
	models.ShowsSet(shows)

	//Calculate time to run
	diff := time.Now().Sub(t1).Seconds()
	fmt.Printf("%s [shows with-shows]/[total]: %d/%d, this took: %fs\n", time.Now().Format("2006-01-02 15:04:05"), showsWithSeasons, len(shows), diff)
	fmt.Printf("%s Shows not processed: %d\n", time.Now().Format("2006-01-02 15:04:05"), notProcessed)
}
