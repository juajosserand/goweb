package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gituhb.com/juajosserand/goweb/internal/product"
)

var (
	errNilService         = errors.New("invalid nil product service")
	errInvalidId          = errors.New("invalid product id")
	errInvalidPrice       = errors.New("invalid product price")
	errCreation           = errors.New("unable to create product")
	errInvalidProductData = errors.New("invalid product data")
)

type ProductHandlers struct {
	svc product.ProductService
}

func NewProductHandlers(s product.ProductService) (ProductHandlers, error) {
	if s == nil {
		return ProductHandlers{}, errNilService
	}

	return ProductHandlers{
		svc: s,
	}, nil
}

type request struct {
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,gte=1"`
	CodeValue   string  `json:"code_value" binding:"required,uppercase,alphanum"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required,gte=0"`
}

func (ph *ProductHandlers) Pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func (ph *ProductHandlers) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"products": ph.svc.All(),
	})
}

func (ph *ProductHandlers) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidId.Error(),
		})
		return
	}

	p, err := ph.svc.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidId.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func (ph *ProductHandlers) Search(ctx *gin.Context) {
	price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidPrice.Error(),
		})
		return
	}

	ps, err := ph.svc.PriceGreaterThan(price)
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

func (ph *ProductHandlers) Create(ctx *gin.Context) {
	var r request

	err := ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidProductData.Error(),
		})
		return
	}

	// create product
	err = ph.svc.Create(
		r.Name,
		r.Quantity,
		r.CodeValue,
		r.IsPublished,
		r.Expiration,
		r.Price,
	)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errCreation.Error(),
		})
	}

	ctx.Status(http.StatusCreated)
}
