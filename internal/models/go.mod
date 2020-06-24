module github.com/egeback/play_media_api/internal/models

go 1.14

replace (
	github.com/egeback/play_media_api/internal/controllers => ../controllers
	github.com/egeback/play_media_api/internal/docs => ../docs
	github.com/egeback/play_media_api/internal/models => ../models
	github.com/egeback/play_media_api/internal/parsers => ../parsers
	github.com/egeback/play_media_api/internal/parsers/svtplay => ../parsers/svtplay
	github.com/egeback/play_media_api/internal/parsers/tv4play => ../parsers/tv4play
	github.com/egeback/play_media_api/internal/utils => ../utils
)

require github.com/egeback/play_media_api/internal/utils v0.0.0-20200622193517-7b39954865db
