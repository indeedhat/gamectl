package api

import (
	"net/http"

	"github.com/indeedhat/command-center/app/config"

	"github.com/gin-gonic/gin"
)

// StartAppController will attempt to start an aplication on the server
func StartAppController(ctx *gin.Context) {
	appKey := ctx.Param("app_key")

	app := config.GepApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := app.Start(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsonSuccess)
}

// StopAppController will attempt to stop an applicaton on the server
func StopAppController(ctx *gin.Context) {
	appKey := ctx.Param("app_key")

	app := config.GepApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := app.Stop(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsonSuccess)
}

// RestartAppController will attempt to first stop and then start the application on the server
//
// I the stop command fails because the app is already stopped then it will just run the start command
func RestartAppController(ctx *gin.Context) {
	appKey := ctx.Param("app_key")

	app := config.GepApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// TODO: Dont really care about the stop for testing, this will need to be rectified once MVP is done
	_ = app.Stop()

	if err := app.Start(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsonSuccess)
}

// GetAppStatusController will run the app status controller and return its json
func GetAppStatusController(ctx *gin.Context) {
	appKey := ctx.Param("app_key")

	app := config.GepApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	status, err := app.Status()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"outcome": true,
		"status":  status,
	})
}
