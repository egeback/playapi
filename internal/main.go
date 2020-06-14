package main

import (
	"fmt"

	"./api"
	"./parsers"
	parser "github.com/egeback/play_media_api/internal/parsers/svtplay"
)

var (
	// BuildVersion updated when running build script
	BuildVersion string = ""
	// BuildTime updated when running build script
	BuildTime string = ""
)

func main() {
	fmt.Println("Running Play Media API version: %s (%s)", BuildVersion, BuildTime)
	parsers := []parsers.ParserInterface{new(parser.SvtPlayParser), new(parser.SvtPlayParser)}
	//parsers := []parsers.ParserInterface{new(tv4play.Parser)}

	api := api.Init(parsers)
	api.Start()
}
