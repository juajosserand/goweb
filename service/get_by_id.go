package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *server) getById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Errors = append(ctx.Errors, ctx.Error(errInvalidId))
		ctx.Status(http.StatusBadRequest)
		return
	}

	p, err := s.repo.GetById(id)
	if err != nil {
		ctx.Errors = append(ctx.Errors, ctx.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, p)
}
