package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/command-center/app/config"
)

func IndexController(ctx *gin.Context) {
	view(ctx, "index", gin.H{
		"apps": config.Apps(),
	})
}
