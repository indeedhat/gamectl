package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/gamectl/internal/config"
)

// ReloadAppConfig
//
// i got sick of having to restart the server and log back in when testing changes to the config
// files so i added this
func ReloadAppConfig(ctx *gin.Context) {
	outcome := config.ReloadAppConfig() == nil

	if outcome {
		ctx.Header("Refresh", "3;url=/")
	}

	view(ctx, "system/reload", gin.H{
		"outcome": outcome,
	})
}
