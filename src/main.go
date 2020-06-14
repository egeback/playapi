package main

import (
	"./api"
	"./parsers"
)

func main() {
	parsers := []parsers.ParserInterface{new(parser.SvtPlayParser), new(parser.SvtPlayParser)}
	//parsers := []parsers.ParserInterface{new(tv4play.Parser)}

	api := api.Init(parsers)
	api.Start()
}
