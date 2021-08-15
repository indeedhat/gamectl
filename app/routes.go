package app

import (
	"encoding/json"
	"html/template"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/indeedhat/command-center/app/controllers"
	"github.com/indeedhat/command-center/app/controllers/api"
	"github.com/indeedhat/command-center/app/middleware"

	"github.com/gin-gonic/gin"
)

var (
	sessionStore memstore.Store
)

// BuildRoutes will setup the gin router for handling web requests along with
// setting up static fs bindings for serving assets
func BuildRoutes() *gin.Engine {
	router := gin.Default()

	setupTemplateFunctions(router)
	setupStatics(router)
	setupSessions(router)

	public := router.Group("/", middleware.IsGuest)
	{
		public.GET("/login", controllers.LoginController)
		public.POST("/login", controllers.LoginController)
	}

	private := router.Group("/", middleware.IsLoggedIn)
	{
		private.GET("/", controllers.IndexController)
		private.GET("/logout", controllers.LogoutController)

		private.GET("/api/apps/:app_key", api.GetAppStatusController)
		private.POST("/api/apps/:app_key/start", api.StartAppController)
		private.POST("/api/apps/:app_key/stop", api.StopAppController)
		private.POST("/api/apps/:app_key/restart", api.RestartAppController)

		private.GET("/ap/apps/:app_key/config/:config_key", api.LoadAppConfig)
		private.POST("/ap/apps/:app_key/config/:config_key", api.SaveAppConfig)
	}

	return router
}

// setupStatics will setup static bindings for asset file serving
// along with assign the templating engines its views directory
func setupStatics(router *gin.Engine) {
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("./views/**/*")
}

// setupSessions for use with gin for user sessions etc
func setupSessions(router *gin.Engine) {
	sessionStore = memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	router.Use(sessions.Sessions("session", sessionStore))
}

func setupTemplateFunctions(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"json": func(data interface{}) string {
			bytes, err := json.Marshal(data)
			if err != nil {
				return ""
			}

			return string(bytes)
		},
	})
}
