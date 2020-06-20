module github.com/egeback/play_media_api/internal

go 1.14

replace (
	github.com/egeback/play_media_api/internal/controllers => ./controllers
	github.com/egeback/play_media_api/internal/docs => ./docs
	github.com/egeback/play_media_api/internal/models => ./models
	github.com/egeback/play_media_api/internal/parsers => ./parsers
	github.com/egeback/play_media_api/internal/parsers/svtplay => ./parsers/svtplay
	github.com/egeback/play_media_api/internal/parsers/tv4play => ./parsers/tv4play
	github.com/egeback/play_media_api/internal/utils => ./utils
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/egeback/play_media_api/internal/controllers v0.0.0-00010101000000-000000000000
	github.com/egeback/play_media_api/internal/docs v0.0.0-00010101000000-000000000000
	github.com/egeback/play_media_api/internal/models v0.0.0-20200620141031-e00657304b36
	github.com/egeback/play_media_api/internal/parsers v0.0.0-20200620141031-e00657304b36
	github.com/egeback/play_media_api/internal/parsers/svtplay v0.0.0-20200615044914-645054921c99
	github.com/egeback/play_media_api/internal/parsers/tv4play v0.0.0-20200615044914-645054921c99
	github.com/egeback/play_media_api/internal/version v0.0.0-20200615181031-ff340633ceaf
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/gosimple/slug v1.9.0 // indirect
	github.com/jasonlvhit/gocron v0.0.0-20200423141508-ab84337f7963
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/liip/sheriff v0.0.0-20190308094614-91aa83a45a3d // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/tomwei7/gin-jsonp v0.0.0-20191103091125-e5236eb5393d
	github.com/urfave/cli/v2 v2.2.0 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200619210111-0f592d2728bb // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
