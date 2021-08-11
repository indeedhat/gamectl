package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/command-center/app/models"
)

// loginInput
type loginInput struct {
	Username string `form:"username" binding:"required"`
	Passwd   string `form:"passwd" binding:"required"`
}

// LoginController displays the logn for and handles the login logic
func LoginController(ctx *gin.Context) {
	var input loginInput
	var errorMessage string

	err := ctx.Bind(&input)
	if err != nil {
		if input.Username != "" && input.Passwd != "" {
			errorMessage = "Bad input"
		}
	} else {
		user := models.LoadUserByLoginDetails(input.Username, input.Passwd)
		if user == nil {
			errorMessage = "Bad login details"
		} else {
			ses := sessions.Default(ctx)
			ses.Set("ID", user.ID)

			ctx.Redirect(http.StatusFound, "/")
			return
		}
	}

	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		"input": input,
		"error": errorMessage,
	})
}

// LogoutController will clear the users session
func LogoutController(ctx *gin.Context) {
	ses := sessions.Default(ctx)

	ses.Delete("ID")

	ctx.Redirect(http.StatusFound, "/login")
}
