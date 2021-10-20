package middleware

import (
	"log"

	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func DebugSessions(store *memstore.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		log.Printf("%#v\n", *store)
		ctx.Next()
	}
}
