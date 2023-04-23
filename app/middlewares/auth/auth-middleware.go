package middleware_auth

import (
	"net/http"
	"simple-hosting/compositor/app/settings"

	"github.com/gin-gonic/gin"
)

var config = settings.GetServiceSettings()

func CheckAuth(ctx *gin.Context) {
	apikey := ctx.GetHeader("X-ApiKey")
	if apikey != config.App.ApiKey {
		ctx.JSON(http.StatusUnauthorized, "You should specify your API key in the X-ApiKey header")
		ctx.Abort()
	}
}
