package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-version"
	"github.com/liip/sheriff"
)

var (
	apiVersion = "1.0.0"
)

// Controller struct
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

//Response used to describe api response
type Response struct {
	Data interface{} `json:"data" groups:"api"`
}

//ErrorResponse used to descrive api error response
type ErrorResponse struct {
	Message  string `json:"message" groups:"api" example:"status bad request"`
	Code     int    `json:"code" groups:"api" example:"400"`
	MoreInfo string `json:"more_info" groups:"api" example:"http://"`
}

func (c *Controller) createJSONResponse(ctx *gin.Context, obj interface{}, groups ...string) {
	data, err := marshal(obj, false, groups...)
	if err != nil {
		c.createErrorResponse(ctx, 500, 100, "Could not marshal response")
	}
	ctx.Data(http.StatusOK, "application/json", data)
}

func (c *Controller) createJSONResponsePretty(ctx *gin.Context, obj interface{}, groups ...string) {
	data, err := marshal(obj, true, groups...)
	if err != nil {
		c.createErrorResponse(ctx, 500, 100, "Could not marshal response")
	}
	ctx.Data(http.StatusOK, "application/json", data)
}

func (c *Controller) createErrorResponse(ctx *gin.Context, statusCode int, code int, message string) {
	ctx.JSON(statusCode, ErrorResponse{
		Message:  message,
		Code:     code,
		MoreInfo: fmt.Sprint("/%s", code),
	})
}

func marshal(data interface{}, prettyPrint bool, groups ...string) ([]byte, error) {
	v1, err := version.NewVersion(apiVersion)
	if err != nil {
		log.Panic(err)
		return []byte{}, err
	}

	groups = append(groups, "api")
	o := &sheriff.Options{
		Groups:     groups,
		ApiVersion: v1,
	}

	dest, err := sheriff.Marshal(o, data)
	if err != nil {
		return nil, err
	}
	if prettyPrint {
		return json.MarshalIndent(dest, "", "  ")
	}
	return json.Marshal(dest)
}
