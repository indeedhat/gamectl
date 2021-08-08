package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexController(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
        "world": "World!",
        "counter": []string{
            "one",
            "two",
            "three",
            "four",
        },
    })
}
