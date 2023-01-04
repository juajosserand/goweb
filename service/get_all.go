package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) getAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"products": s.repo.All(),
	})
}
