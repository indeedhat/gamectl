package api

import (
	"io/ioutil"
	"net/http"

	"github.com/indeedhat/command-center/app/config"

	"github.com/gin-gonic/gin"
)

type saveAppConfigInput struct {
	Data string `form:"data" bindings:"required"`
}

// LoadAppConfig will return the contents of a specific app config file
func LoadAppConfig(ctx *gin.Context) {
	appKey := ctx.Param("app_key")
	configKey := ctx.Param("config_key")

	app := config.GepApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	info, ok := app.Files[configKey]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	contents, err := ioutil.ReadFile(info.Path)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"outcome": true,
		"file":    string(contents),
	})
}

// SaveAppConfig will attempt to update the specific app config
func SaveAppConfig(ctx *gin.Context) {
	var input saveAppConfigInput
	appKey := ctx.Param("app_key")
	configKey := ctx.Param("config_key")

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	app := config.GepApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	info, ok := app.Files[configKey]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := ioutil.WriteFile(info.Path, []byte(input.Data), 0644); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsonSuccess)
}
