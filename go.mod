module github.com/egeback/play_media_api

go 1.14

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/egeback/play_media_api/internal/controllers v0.0.0-00010101000000-000000000000 // indirect
	github.com/egeback/play_media_api/internal/docs v0.0.0-00010101000000-000000000000 // indirect
	github.com/egeback/play_media_api/internal/models v0.0.0-20200615181031-ff340633ceaf // indirect
	github.com/egeback/play_media_api/internal/version v0.0.0-20200615181031-ff340633ceaf // indirect
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/jasonlvhit/gocron v0.0.0-20200423141508-ab84337f7963 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14 // indirect
	github.com/swaggo/gin-swagger v1.2.0 // indirect
	github.com/tomwei7/gin-jsonp v0.0.0-20191103091125-e5236eb5393d // indirect
)

replace (
	github.com/egeback/play_media_api/internal/controllers => ./internal/controllers
	github.com/egeback/play_media_api/internal/docs => ./internal/docs
	github.com/egeback/play_media_api/internal/models => ./internal/models
	github.com/egeback/play_media_api/internal/parsers => ./internal/parsers
	github.com/egeback/play_media_api/internal/parsers/svtplay => ./internal/parsers/svtplay
	github.com/egeback/play_media_api/internal/parsers/tv4play => ./internal/parsers/tv4play
)
