package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/command-center/app/models"
)

// IsLoggedIn will attempt to load the user from session
//
// User will be redirected to login page if not
func IsLoggedIn(ctx *gin.Context) {
	ses := sessions.Default(ctx)

	userId := ses.Get("ID")
	if userId == "0" || userId == nil {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.AbortWithStatus(http.StatusFound)
		return
	}

	user := models.FindUser(userId.(string))
	if user == nil {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.AbortWithStatus(http.StatusFound)
		return
	}

	ctx.Set("user", user)

	ctx.Next()
}

// IsRoot checks that the user is actually a root user
func IsRoot(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.AbortWithStatus(http.StatusFound)
		return
	}

	userObject, ok := user.(*models.User)
	if !ok {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.AbortWithStatus(http.StatusFound)
		return
	}

	if !userObject.Root {
		ctx.Redirect(http.StatusFound, "/users")
		ctx.AbortWithStatus(http.StatusFound)
		return
	}

	ctx.Next()
}

// IsGues check that the user is not logged in
//
// User will be redirected to the panel if not
func IsGuest(ctx *gin.Context) {
	ses := sessions.Default(ctx)

	userId := ses.Get("ID")
	if userId != 0 && userId != nil {
		user := models.FindUser(userId.(string))
		if user != nil {
			ctx.Redirect(http.StatusFound, "/")
			ctx.AbortWithStatus(http.StatusFound)
			return
		}

		ses.Delete("ID")
		ses.Save()
	}

	ctx.Next()
}
