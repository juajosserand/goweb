package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "pong"})
}
