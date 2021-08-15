package controllers

import (
	"fmt"
	"net/http"

	"github.com/indeedhat/command-center/app/models"

	"github.com/gin-gonic/gin"
)

// View is a helper for more cleanly displaying vies
func view(ctx *gin.Context, template string, data gin.H) {
	user, exists := ctx.Get("user")
	if exists {
		fmt.Println("user exits")
		if userObject, ok := user.(*models.User); ok {
			fmt.Println("user added")
			data["user"] = userObject
		}
	}

	ctx.HTML(http.StatusOK, template, data)
}
