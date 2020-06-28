package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping return pong
// @Summary ping common
// @Description do ping
// @Tags common
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /ping [get]
func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "pong")
	return
}
