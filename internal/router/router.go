package router

import (
	"encoding/json"
	"html/template"
	"os"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/indeedhat/gamectl/internal/controllers"
	"github.com/indeedhat/gamectl/internal/controllers/api"
	"github.com/indeedhat/gamectl/internal/middleware"

	"github.com/gin-gonic/gin"
)

var (
	sessionStore memstore.Store
)

// BuildRoutes will setup the gin router for handling web requests along with
// setting up static fs bindings for serving assets
func BuildRoutes() *gin.Engine {
	router := gin.Default()

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
		private.GET("/users", controllers.ListUsersController)
		private.GET("/users/passwd", controllers.UpdatePasswordController)
		private.POST("/users/passwd", controllers.UpdatePasswordController)

		private.GET("/api/apps/:app_key", api.GetAppStatusController)
		private.POST("/api/apps/:app_key/start", api.StartAppController)
		private.POST("/api/apps/:app_key/stop", api.StopAppController)
		private.POST("/api/apps/:app_key/restart", api.RestartAppController)
		private.GET("/api/apps/:app_key/download", api.DownloadAppWorldController)

		private.GET("/api/apps/:app_key/config/:config_key", api.LoadAppConfig)
		private.POST("/api/apps/:app_key/config/:config_key", api.SaveAppConfig)

		private.GET("/api/apps/:app_key/logs/:log_key", api.LogStreamController)

		private.GET("/api/performance", api.SystemPerformanceStreamController)
	}

	rootAdmin := router.Group("/", middleware.IsLoggedIn, middleware.IsRoot)
	{
		rootAdmin.GET("/users/:user_id", controllers.UpdateUserController)
		rootAdmin.POST("/users/:user_id", controllers.UpdateUserController)

		rootAdmin.GET("/users/create", controllers.CreateUserController)
		rootAdmin.POST("/users/create", controllers.CreateUserController)

		rootAdmin.GET("/system/reload", controllers.ReloadAppConfig)

		private.GET("/ws/apps/:app_key/tty", api.TtySocketController)
	}

	return router
}

// setupStatics will setup static bindings for asset file serving
// along with assign the templating engines its views directory
func setupStatics(router *gin.Engine) {
	static := router.Group("/", gzip.Gzip(gzip.DefaultCompression))
	static.Static("/assets", "./web/assets")

	viewsConfig := goview.DefaultConfig
	viewsConfig.Root = "./web/views"
	viewsConfig.DisableCache = true
	viewsConfig.Funcs = template.FuncMap{
		"json": func(data interface{}) string {
			bytes, err := json.Marshal(data)
			if err != nil {
				return ""
			}

			return string(bytes)
		},
	}

	router.HTMLRender = ginview.New(viewsConfig)
}

// setupSessions for use with gin for user sessions etc
func setupSessions(router *gin.Engine) {
	sessionStore = memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	router.Use(sessions.Sessions("session", sessionStore))
}
