package controllers

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/egeback/playapi/internal/models"
	"github.com/egeback/playapi/internal/utils"
	"github.com/gin-gonic/gin"
)

//ListLatestEpisodes function returns the latest episodes from API
// @Summary List latest episodes
// @Description get episodes
// @Tags episodes
// @Accept json
// @Produce json
// @Param prettyPrint query string false "pretty print show" Format(bool)
// @Success 200 {array} models.Episode
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /episodes/latest [get]
func (c *Controller) ListLatestEpisodes(ctx *gin.Context) {
	prettyPrint := utils.GetBoolValueFromString(ctx.DefaultQuery("prettyPrint", ""), false)

	limit := utils.GetIntValueFromString(ctx.DefaultQuery("limit", ""), *defaultLimit)
	offset := utils.GetIntValueFromString(ctx.DefaultQuery("offset", ""), 0)

	//q, err := extractQueryParameter(ctx.DefaultQuery("q", ""))
	var err error

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

	episodes, err := models.GetLatestEpisodes()
	if err != nil {
		log.Println(err)
		c.createErrorResponse(ctx, 500, 100, "Could not fetch episodes")
	}

	size := len(episodes)

	if *offset > size {
		episodes = make([]models.Episode, 0, 0)
	} else if *offset+*limit > size {
		episodes = episodes[*offset:]
	} else {
		episodes = episodes[*offset : *offset+*limit]
	}

	//Create paged response struct
	response := PagedResponse{
		Limit:   utils.Min(*limit, size),
		Size:    size,
		Start:   *offset,
		Results: episodes,
	}

	if *prettyPrint {
		j, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			c.createErrorResponse(ctx, 500, 100, "Could not marshal response")
			return
		}
		ctx.Data(200, "application/json", j)
	} else {
		ctx.JSON(200, response)
		return
	}
}

//ListEpisodes function returns the latest episodes from API
// @Summary List all episodes
// @Description get episodes
// @Tags episodes
// @Accept json
// @Produce json
// @Param prettyPrint query string false "pretty print show" Format(bool)
// @Success 200 {array} models.Episode
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /episodes [get]
func (c *Controller) ListEpisodes(ctx *gin.Context) {
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

	episodes, err := models.GetEpisodes(q...)
	if err != nil {
		log.Println(err)
		c.createErrorResponse(ctx, 500, 100, "Could not fetch episodes")
	}

	size := len(episodes)

	if *offset > size {
		episodes = make([]models.Episode, 0, 0)
	} else if *offset+*limit > size {
		episodes = episodes[*offset:]
	} else {
		episodes = episodes[*offset : *offset+*limit]
	}

	//Create paged response struct
	response := PagedResponse{
		Limit:   utils.Min(*limit, size),
		Size:    size,
		Start:   *offset,
		Results: episodes,
	}

	if *prettyPrint {
		j, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			c.createErrorResponse(ctx, 500, 100, "Could not marshal response")
			return
		}
		ctx.Data(200, "application/json", j)
	} else {
		ctx.JSON(200, response)
		return
	}
}
