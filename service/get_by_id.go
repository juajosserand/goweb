package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *server) getById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidId.Error(),
		})
		return
	}

	p, err := s.repo.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidId.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, p)
}
