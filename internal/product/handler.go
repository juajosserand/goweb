package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	svc ProductService
}

func NewHandler(mux *gin.Engine, s ProductService) {
	ph := &productHandler{
		svc: s,
	}

	mux.GET("/ping", ph.Pong)

	productsMux := mux.Group("/products")
	productsMux.GET("/", ph.GetAll)
	productsMux.GET("/:id", ph.GetById)
	productsMux.GET("/search", ph.Search)
	productsMux.POST("/", ph.Create)
	productsMux.PUT("/:id", ph.Update)
	productsMux.PATCH("/:id", ph.UpdateName)
	productsMux.DELETE("/:id", ph.Delete)
}

type request struct {
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,gte=1"`
	CodeValue   string  `json:"code_value" binding:"required,uppercase,alphanum"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required,gte=0"`
}

type patchRequest struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (ph *productHandler) Pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func (ph *productHandler) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"products": ph.svc.All(),
	})
}

func (ph *productHandler) GetById(ctx *gin.Context) {
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

func (ph *productHandler) Search(ctx *gin.Context) {
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

func (ph *productHandler) Create(ctx *gin.Context) {
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
		return
	}

	ctx.Status(http.StatusCreated)
}

func (ph *productHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidId.Error(),
		})
		return
	}

	var r request

	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidProductData.Error(),
		})
		return
	}

	err = ph.svc.Update(
		id,
		r.Name,
		r.Quantity,
		r.CodeValue,
		r.IsPublished,
		r.Expiration,
		r.Price,
	)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errInvalidProductData.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *productHandler) UpdateName(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidId.Error(),
		})
		return
	}

	var r patchRequest

	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidProductData.Error(),
		})
		return
	}

	err = ph.svc.UpdateName(
		id,
		r.Name,
	)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errInvalidProductData.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *productHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidId.Error(),
		})
		return
	}

	err = ph.svc.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": errDeletion.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
