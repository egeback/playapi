module github.com/egeback/play_media_api

go 1.14

require (
	github.com/egeback/play_media_api/internal/models v0.0.0-20200615181031-ff340633ceaf // indirect
	github.com/egeback/play_media_api/internal/version v0.0.0-20200615181031-ff340633ceaf // indirect
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/jasonlvhit/gocron v0.0.0-20200423141508-ab84337f7963 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14 // indirect
	github.com/swaggo/gin-swagger v1.2.0 // indirect
)

replace (
	github.com/egeback/play_media_api/internal/api => internal/api
	github.com/egeback/play_media_api/internal/controllers => internal/controllers
	github.com/egeback/play_media_api/internal/docs => internal/docs
	github.com/egeback/play_media_api/internal/models => internal/models
	github.com/egeback/play_media_api/internal/parsers => internal/parsers
	github.com/egeback/play_media_api/internal/parsers/svtplay => internal/parsers/svtplay
	github.com/egeback/play_media_api/internal/parsers/tv4play => internal/parsers/tv4play
)
