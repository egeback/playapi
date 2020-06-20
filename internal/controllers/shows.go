package controllers

import (
	"log"
	"strconv"
	"strings"

	"github.com/egeback/play_media_api/internal/models"
	"github.com/gin-gonic/gin"
)

//ListShows function to return shows from API
// @Summary List shows
// @Description get shows
// @Tags shows
// @Accept  json
// @Produce  json
// @Param prettyPrint string false "pretty print show" Format(bool)
// @Success 200 {array} models.Show
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /shows [get]
func (c *Controller) ListShows(ctx *gin.Context) {
	seasons := ctx.DefaultQuery("seasons", "false")
	displaySeasons, err := strconv.ParseBool(seasons)
	pp := ctx.DefaultQuery("prettyPrint", "false")
	prettyPrint, err2 := strconv.ParseBool(pp)

	if err != nil || err2 != nil {
		log.Panic(err)
		log.Panic(err2)
		ctx.JSON(200, gin.H{
			"data":       "coudld not parse quyery parameters",
			"status":     "error",
			"statusCode": 400,
			"data_type":  "String",
		})
	}

	shows, err := models.ShowsAll("")
	if err != nil {
		log.Panic(err)
		c.createErrorResponse(ctx, 500, 100, "Could not fetch shows")
	}

	if !displaySeasons {
		if prettyPrint {
			c.createJSONResponsePretty(ctx, shows)
		} else {
			c.createJSONResponse(ctx, shows)
		}
	} else {
		if prettyPrint {
			c.createJSONResponsePretty(ctx, shows, "seasons")
		} else {
			c.createJSONResponse(ctx, shows, "seasons")
		}
	}
}

//ShowShow func to retrun a specific show
// @Summary Show an show
// @Description get show by slug
// @Tags shows
// @Accept json
// @Produce json
// @Param slug path string true "Show ID"
// @Param prettyPrint string false "pretty print show" Format(bool)
// @Success 200 {object} models.Show
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /show/{slug} [get]
func (c *Controller) ShowShow(ctx *gin.Context) {
	slug := ctx.Param("slug")
	prettyPrint := ctx.DefaultQuery("prettyPrint", "false")
	var show *models.Show

	shows, err := models.ShowsAll(slug)
	if err != nil {
		c.createErrorResponse(ctx, 500, 100, "Could not fetch shows")
		return
	}

	if len(shows) == 0 {
		c.createErrorResponse(ctx, 404, 101, "Show not found")
		return
	}

	show = &shows[0]

	if strings.ToLower(prettyPrint) == "false" {
		c.createJSONResponse(ctx, show, "seasons")
		return
	}
	c.createJSONResponsePretty(ctx, show, "seasons")
	return
}
