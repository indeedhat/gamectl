package api

import (
	"io"
	"net/http"
	"time"

	"github.com/indeedhat/gamectl/internal/config"

	"github.com/gin-gonic/gin"
)

// LogStreamController setups up persistant connection that it uses to send log updates to the client
func LogStreamController(ctx *gin.Context) {
	appKey := ctx.Param("app_key")
	logKey := ctx.Param("log_key")

	app := config.GetApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	serverLog, ok := app.Logs[logKey]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	done := make(chan bool)
	logUpdated, err := serverLog.Watch(done)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	keepAlive := time.NewTicker(time.Second * 55)
	clientDisconnected := ctx.Writer.CloseNotify()
	ctx.Stream(func(writer io.Writer) bool {
		select {
		case <-clientDisconnected:
			done <- true
			return false

		case <-keepAlive.C:
			ctx.SSEvent("keep-alive", true)
			return true

		case data := <-logUpdated:
			ctx.SSEvent("message", data)
			return true
		}
	})
}
