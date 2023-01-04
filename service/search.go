package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *server) search(ctx *gin.Context) {
	price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidPrice.Error(),
		})
		return
	}

	ps, err := s.repo.PriceGreaterThan(price)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidPrice.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": ps,
	})
}
