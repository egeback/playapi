package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/egeback/playapi/internal/models"
	"github.com/gin-gonic/gin"
)

//ListShows function to return shows from API
// @Summary List shows
// @Description get shows
// @Tags shows
// @Accept json
// @Produce json
// @Param prettyPrint query string false "pretty print show" Format(bool)
// @Param showSeasons query string false "show seasons" Format(bool)
// @Success 200 {array} models.Show
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /shows [get]
func (c *Controller) ListShows(ctx *gin.Context) {
	seasons := ctx.DefaultQuery("showSeasons", "false")
	displaySeasons, err := strconv.ParseBool(seasons)

	pp := ctx.DefaultQuery("prettyPrint", "false")
	prettyPrint, err2 := strconv.ParseBool(pp)

	q := extractQueryParameter(ctx.DefaultQuery("q", ""))

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

	shows, err := models.ShowsAll(q...)
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
// @Param prettyPrint query string false "pretty print show" Format(bool)
// @Success 200 {object} models.Show
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /shows/{slug} [get]
func (c *Controller) ShowShow(ctx *gin.Context) {
	slug := ctx.Param("slug")
	prettyPrint := ctx.DefaultQuery("prettyPrint", "false")
	var show *models.Show

	shows, err := models.ShowsAll(models.QueryItem{Field: "slug", Value: slug})
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

func extractQueryParameter(str string) []models.QueryItem {
	queryItems := make([]models.QueryItem, 0, 0)
	if len(str) == 0 {
		return queryItems
	}

	err := json.Unmarshal([]byte(str), &queryItems)
	if err != nil {
		log.Panic(err)
		return queryItems
	}

	//a := strings.Split(str, ";")
	// for _, s := range a {
	// 	parts := strings.Split(s, ":")
	// 	if len(parts) == 2 {
	// 		if strings.Index(s, ",") >= 0 {
	// 			queryItems = append(queryItems, models.QueryItem{Filter: parts[0], Value: strings.Split(parts[1], ",")})
	// 		} else {
	// 			queryItems = append(queryItems, models.QueryItem{Filter: parts[0], Value: parts[1]})
	// 		}
	// 	}
	// }

	return queryItems

}
