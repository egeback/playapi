module github.com/egeback/play_media_api/internal/controllers

go 1.14

replace github.com/egeback/play_media_api/internal/models => ../models

require (
	github.com/egeback/play_media_api/internal v0.0.0-20200615181031-ff340633ceaf // indirect
	github.com/egeback/play_media_api/internal/models v0.0.0-20200615181031-ff340633ceaf
	github.com/egeback/play_media_api/internal/version v0.0.0-20200615181031-ff340633ceaf // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/hashicorp/go-version v1.2.1
	github.com/liip/sheriff v0.0.0-20190308094614-91aa83a45a3d
)
