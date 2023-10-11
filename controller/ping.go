package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct{}

func (c *PingController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
