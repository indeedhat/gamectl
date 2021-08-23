package api

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/indeedhat/gamectl/app/performance"
)

// LogStreamController setups up persistant connection that it uses to send log updates to the client
func SystemPerformanceStreamController(ctx *gin.Context) {
	clientDisconnected := ctx.Writer.CloseNotify()
	ctx.Stream(func(writer io.Writer) bool {
		select {
		case <-clientDisconnected:
			return false

		case <-time.After(time.Second * performance.PollingInterval):
			monitor := performance.GetMonitor()
			if monitor == nil {
				return true
			}

			data := monitor.Read()
			json, _ := data.Json()

			ctx.SSEvent("message", json)
			return true
		}
	})
}
