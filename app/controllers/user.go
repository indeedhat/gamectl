package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/indeedhat/command-center/app/models"
)

// loginInput form input for logging in
type loginInput struct {
	Username string `form:"username" binding:"required"`
	Passwd   string `form:"passwd" binding:"required"`
}

// updateUserInput form input for updatig a user
type updateUserInput struct {
	Username string `form:"username" binding:"required"`
	Passwd   string `form:"passwd"`
}

// createUserInput form input for creating a new user
type createUserInput struct {
	Username string `form:"username" binding:"required"`
	Passwd   string `form:"passwd" binding:"required"`
}

// updatePasswordInput form input for a user to update their own password
type updatePasswordInput struct {
	Passwd string `form:"passwd" binding:"required"`
	Verify string `form:"verify" binding:"required"`
}

// LoginController displays the logn for and handles the login logic
func LoginController(ctx *gin.Context) {
	var input loginInput
	var errorMessage string

	err := ctx.Bind(&input)
	if err != nil {
		if input.Username != "" || input.Passwd != "" {
			errorMessage = "Bad input"
		}
	} else {
		user := models.LoadUserByLoginDetails(input.Username, input.Passwd)
		if user == nil {
			errorMessage = "Bad login details"
		} else {
			ses := sessions.Default(ctx)
			ses.Set("ID", strconv.Itoa(int(user.ID)))
			ses.Save()

			ctx.Redirect(http.StatusFound, "/")
			ctx.AbortWithStatus(http.StatusFound)
			return
		}
	}

	view(ctx, "login", gin.H{
		"input": input,
		"error": errorMessage,
	})
}

// LogoutController will clear the users session
func LogoutController(ctx *gin.Context) {
	ses := sessions.Default(ctx)

	ses.Delete("ID")
	ses.Save()

	ctx.Redirect(http.StatusFound, "/login")
}

// ListUsersController displays the full user list
func ListUsersController(ctx *gin.Context) {
	users := models.ListUsers()

	if users == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	view(ctx, "users/index", gin.H{
		"users": users,
	})
}

// UpdateUserController will display the update form and handle the submit
func UpdateUserController(ctx *gin.Context) {
	var input updateUserInput
	var errorString string

	userId := ctx.Param("user_id")

	user := models.FindUser(userId)
	if user == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	err := ctx.Bind(&input)
	if err != nil {
		if input.Username != "" && input.Passwd != "" {
			errorString = "Bad Input"
		}
	} else {
		err := models.UpdateUser(user, input.Username, input.Passwd)
		if err != nil {
			errorString = "Update Failed"
		} else {
			ctx.Redirect(http.StatusFound, "/users")
			ctx.AbortWithStatus(http.StatusFound)
			return
		}
	}

	view(ctx, "users/update", gin.H{
		"subject":     user,
		"errorString": errorString,
		"input":       input,
	})

}

// CreateUserController will show the create user form and attempt to actually create the user
func CreateUserController(ctx *gin.Context) {
	var input createUserInput
	var errorString string

	err := ctx.Bind(&input)
	if err != nil {
		if input.Username != "" && input.Passwd != "" {
			errorString = "Bad Input"
		}
	} else {
		user := models.CreateUser(input.Username, input.Passwd)
		if user == nil {
			errorString = "Create Failed"
		} else {
			ctx.Redirect(http.StatusFound, "/users")
			ctx.AbortWithStatus(http.StatusFound)
			return
		}
	}

	view(ctx, "users/create", gin.H{
		"errorString": errorString,
		"input":       input,
	})

}

// UpdatePasswordController will let the user update their password
func UpdatePasswordController(ctx *gin.Context) {
	var input updatePasswordInput
	var errorString string

	userId := ctx.Param("user_id")

	user := models.FindUser(userId)
	if user == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	err := ctx.Bind(&input)
	if err != nil {
		if input.Verify != "" && input.Passwd != "" {
			errorString = "Bad Input"
		}
	} else if input.Verify != input.Passwd {
		errorString = "Passwords do not match"
	} else {
		err := models.UpdateUser(user, user.Name, input.Passwd)
		if err != nil {
			errorString = "Update Failed"
		} else {
			ctx.Redirect(http.StatusFound, "/users")
			ctx.AbortWithStatus(http.StatusFound)
			return
		}
	}

	view(ctx, "users/password", gin.H{
		"subject":     user,
		"errorString": errorString,
		"input":       input,
	})

}
