package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gituhb.com/juajosserand/goweb/model"
)

func (s *server) create(ctx *gin.Context) {
	var p model.Product

	err := ctx.ShouldBindJSON(&p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// code value validation
	for _, product := range s.repo.All() {
		if product.CodeValue == p.CodeValue {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": errDuplicatedCodeValue.Error(),
			})
			return
		}
	}

	// expiration validation
	expDate, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if expDate.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidExpirationDate.Error(),
		})
		return
	}

	var dateStr string

	if expDate.Day() < 10 {
		dateStr += fmt.Sprintf("0%d/", expDate.Day())
	} else {
		dateStr += fmt.Sprintf("%d/", expDate.Day())
	}

	if expDate.Month() < 10 {
		dateStr += fmt.Sprintf("0%d/", expDate.Month())
	} else {
		dateStr += fmt.Sprintf("%d/", expDate.Month())
	}

	dateStr += fmt.Sprintf("%d", expDate.Year())

	p.Expiration = dateStr

	// create product
	s.repo.Create(p)
	ctx.Status(http.StatusCreated)
}
