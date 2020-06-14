package main

import (
	"github.com/egeback/stream_media_api/src/api"
	"github.com/egeback/stream_media_api/src/parsers"

	"github.com/egeback/stream_media_api/src/parsers/svtplay"
	"github.com/egeback/stream_media_api/src/parsers/tv4play"
)

func main() {
	parsers := []parsers.ParserInterface{new(svtplay.Parser), new(tv4play.Parser)}
	//parsers := []parsers.ParserInterface{new(tv4play.Parser)}

	api := api.Init(parsers)
	api.Start()
}
