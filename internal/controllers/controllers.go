package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/indeedhat/gamectl/internal/models"
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
