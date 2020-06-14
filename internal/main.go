package main

import (
	"fmt"

	"github.com/egeback/play_media_api/internal/api"
	"github.com/egeback/play_media_api/internal/parsers"
	"github.com/egeback/play_media_api/internal/parsers/svtplay"
	"github.com/egeback/play_media_api/internal/parsers/tv4play"
)

var (
	// BuildVersion updated when running build script
	BuildVersion string = ""
	// BuildTime updated when running build script
	BuildTime string = ""
)

func main() {
	fmt.Println("Running Play Media API version: %s (%s)", BuildVersion, BuildTime)
	parsers := []parsers.ParserInterface{new(svtplay.Parser), new(tv4play.Parser)}
	//parsers := []parsers.ParserInterface{new(tv4play.Parser)}

	api := api.Init(parsers)
	api.Start()
}
