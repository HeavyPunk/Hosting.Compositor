package middleware_auth

import (
	"net/http"
	"simple-hosting/compositor/app/settings"
	file_settings_provider "simple-hosting/go-commons/settings/file-settings-provider"

	"github.com/gin-gonic/gin"
)

var config = file_settings_provider.GetSetting[settings.ServiceSettings]("settings.yml")

func CheckAuth(ctx *gin.Context) {
	apikey := ctx.GetHeader("X-ApiKey")
	if apikey != config.App.ApiKey {
		ctx.JSON(http.StatusUnauthorized, "You should specify your API key in the X-ApiKey header")
		ctx.Abort()
	}
}
