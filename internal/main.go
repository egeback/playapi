package main

import (
	"fmt"

	"github.com/egeback/play_media_api/internal/api"
	"github.com/egeback/play_media_api/internal/parsers"
	"github.com/egeback/play_media_api/internal/parsers/svtplay"
	"github.com/egeback/play_media_api/internal/parsers/tv4play"
	"github.com/egeback/play_media_api/internal/version"
)

func main() {
	fmt.Printf("Running Play Media API version: %s (%s)\n", version.BuildVersion, version.BuildTime)
	parsers := []parsers.ParserInterface{new(svtplay.Parser), new(tv4play.Parser)}
	//parsers := []parsers.ParserInterface{new(tv4play.Parser)}

	api := api.Init(parsers)
	api.Start()
}
