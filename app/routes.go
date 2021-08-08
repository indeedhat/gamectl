package app

import (
	"github.com/indeedhat/command-center/app/controllers"

	"github.com/gin-gonic/gin"
)

// BuildRoutes will setup the gin router for handling web requests along with
// setting up static fs bindings for serving assets
func BuildRoutes() *gin.Engine {
	router := gin.Default()

	setupStatics(router)

	router.GET("/", controllers.IndexController)

	return router
}

// setupStatics will setup static bindings for asset file serving
// along with assign the templating engines its views directory
func setupStatics(router *gin.Engine) {
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("./views/*")
}
