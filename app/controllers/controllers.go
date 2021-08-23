package controllers

import (
	"net/http"

	"github.com/indeedhat/gamectl/app/models"

	"github.com/gin-gonic/gin"
)

// View is a helper for more cleanly displaying vies
func view(ctx *gin.Context, template string, data gin.H) {
	user, exists := ctx.Get("user")
	if exists {
		if userObject, ok := user.(*models.User); ok {
			data["user"] = userObject
		}
	}

	if _, ok := data["Title"]; !ok {
		data["Title"] = "Command Center"
	}

	ctx.HTML(http.StatusOK, template, data)
}
