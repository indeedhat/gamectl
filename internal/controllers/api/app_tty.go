package api

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/gamectl/internal/config"
)

const (
	writeWait        = 10 * time.Second
	maxMessageSize   = 8192
	pongWait         = 60 * time.Second
	pingPeriod       = (pongWait * 9) / 10
	closeGracePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{}

func TtySocketController(ctx *gin.Context) {
	appKey := ctx.Param("app_key")
	ttyCtx, cancel := context.WithCancel(context.Background())

	app := config.GetApp(appKey)
	if app == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if app.Tty.Command.Command == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	cmd, err := app.Tty.Command.Process()
	if err != nil {
		log.Println("cmd:", err)
		return
	}

	tty, err := pty.Start(cmd)
	if err != nil {
		log.Println("pty:", err)
		return
	}

	go writePump(cancel, ws, tty)
	go ping(ttyCtx, ws)
	readPump(cancel, ws, tty)

	<-ttyCtx.Done()
	log.Print("socket closed")
}

func readPump(cancel context.CancelFunc, ws *websocket.Conn, writer io.Writer) {
	defer ws.Close()
	defer cancel()

	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	log.Print("reading from socket")
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Print(err)
			break
		}

		if _, err := writer.Write(message); err != nil {
			log.Print(err)
			break
		}
	}
}

func writePump(cancel context.CancelFunc, ws *websocket.Conn, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, scanner.Bytes()); err != nil {
			ws.Close()
			break
		}
	}

	if scanner.Err() != nil {
		log.Println("scan:", scanner.Err())
	}

	cancel()

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ctx context.Context, ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				if err.Error() == "websocket: close sent" {
					return
				}
				log.Println("ping:", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
