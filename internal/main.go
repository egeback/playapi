package main

import (
	"fmt"

	playmediaapi "github.com/egeback/play_media_api/internal"
)

var (
	// BuildVersion updated when running build script
	BuildVersion string = ""
	// BuildTime updated when running build script
	BuildTime string = ""
)

func main() {
	fmt.Println("Running Play Media API version: %s (%s)", BuildVersion, BuildTime)
	parsers := []playmediaapi.ParserInterface{new(playmediaapi.SvtPlayParser), new(playmediaapi.SvtPlayParser)}
	//parsers := []parsers.ParserInterface{new(tv4play.Parser)}

	api := api.Init(parsers)
	api.Start()
}
