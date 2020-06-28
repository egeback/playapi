package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/egeback/playapi/internal/models"
	"github.com/egeback/playapi/internal/utils"
	"github.com/gin-gonic/gin"
)

var defaultLimitInt = 100
var defaultLimit = &defaultLimitInt

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
	displaySeasons := utils.GetBoolValueFromString(ctx.DefaultQuery("showSeasons", ""), false)
	prettyPrint := utils.GetBoolValueFromString(ctx.DefaultQuery("prettyPrint", ""), false)

	limit := utils.GetIntValueFromString(ctx.DefaultQuery("limit", ""), *defaultLimit)
	offset := utils.GetIntValueFromString(ctx.DefaultQuery("offset", ""), 0)

	q, err := extractQueryParameter(ctx.DefaultQuery("q", ""))

	if *limit < 0 || *offset < 0 {
		if err == nil {
			err = errors.New("Limit or Offset query parameters < 0")
		}
	}

	if err != nil {
		ctx.JSON(200, gin.H{
			"data":       err.Error(),
			"status":     "error",
			"statusCode": 400,
			"data_type":  "String",
		})
		return
	}

	shows, err := models.ShowsAll(q...)
	if err != nil {
		log.Println(err)
		c.createErrorResponse(ctx, 500, 100, "Could not fetch shows")
	}

	// schema := ""
	// if ctx.Request.TLS == nil {
	// 	schema = "http"
	// } else {
	// 	schema =
	// 		"https"
	// }

	size := len(shows)

	if *offset > size {
		shows = make([]models.Show, 0, 0)
	} else if *offset+*limit > size {
		shows = shows[*offset:]
	} else {
		shows = shows[*offset : *offset+*limit]
	}

	response := PagedResponse{
		Limit:   *limit,
		Size:    size,
		Start:   *offset,
		Results: shows,
	}

	if !*displaySeasons {
		if *prettyPrint {
			c.createJSONResponsePretty(ctx, response)
		} else {
			c.createJSONResponse(ctx, response)
		}
	} else {
		if *prettyPrint {
			c.createJSONResponsePretty(ctx, response, "seasons")
		} else {
			c.createJSONResponse(ctx, response, "seasons")
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

func extractQueryParameter(str string) ([]models.QueryItem, error) {
	queryItems := make([]models.QueryItem, 0, 0)
	if len(str) == 0 {
		return queryItems, nil
	}

	err := json.Unmarshal([]byte(str), &queryItems)
	if err != nil {
		log.Println(err)
		return nil, err
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

	return queryItems, nil

}
