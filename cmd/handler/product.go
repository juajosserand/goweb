package handler

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gituhb.com/juajosserand/goweb/internal/product"
	"gituhb.com/juajosserand/goweb/pkg/web"
)

type productHandler struct {
	svc product.ProductService
}

func NewProductHandler(mux *gin.Engine, s product.ProductService) {
	ph := &productHandler{
		svc: s,
	}

	mux.GET("/ping", ph.Pong)

	productsMux := mux.Group("/products")
	productsMux.GET("/", ph.GetAll)
	productsMux.GET("/:id", ph.GetById)
	productsMux.GET("/search", ph.Search)

	productsMux.Use(auth)

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

// Authentication middleware
func auth(ctx *gin.Context) {
	if ctx.GetHeader("token") != os.Getenv("TOKEN") {
		ctx.JSON(http.StatusUnauthorized, web.ErrResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			product.ErrInvalidToken.Error(),
		))
		return
	}

	ctx.Next()
}

func (ph *productHandler) Pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func (ph *productHandler) GetAll(ctx *gin.Context) {
	ps, err := ph.svc.All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			err.Error(),
		))
		return
	}

	ctx.JSON(http.StatusOK, web.Response(ps))
}

func (ph *productHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidId.Error(),
		))
		return
	}

	p, err := ph.svc.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, web.ErrResponse(
			http.StatusNotFound,
			"Not Found",
			product.ErrInvalidId.Error(),
		))
		return
	}

	ctx.JSON(http.StatusOK, web.Response(p))
}

func (ph *productHandler) Search(ctx *gin.Context) {
	price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidPrice.Error(),
		))
		return
	}

	ps, err := ph.svc.PriceGreaterThan(price)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidPrice.Error(),
		))
		return
	}

	ctx.JSON(http.StatusOK, web.Response(ps))
}

func (ph *productHandler) Create(ctx *gin.Context) {
	var r request

	err := ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidProductData.Error(),
		))
		return
	}

	err = ph.svc.Create(
		r.Name,
		r.Quantity,
		r.CodeValue,
		r.IsPublished,
		r.Expiration,
		r.Price,
	)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, web.ErrResponse(
			http.StatusUnprocessableEntity,
			"Unprocessable Entity",
			product.ErrCreation.Error(),
		))
		return
	}

	ctx.Status(http.StatusCreated)
}

func (ph *productHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidId.Error(),
		))
		return
	}

	var r request

	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidProductData.Error(),
		))
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
		ctx.JSON(http.StatusUnprocessableEntity, web.ErrResponse(
			http.StatusUnprocessableEntity,
			"Unprocessable Entity",
			product.ErrInvalidProductData.Error(),
		))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *productHandler) UpdateName(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidId.Error(),
		))
		return
	}

	type request struct {
		Name        string  `json:"name"`
		Quantity    int     `json:"quantity"`
		CodeValue   string  `json:"code_value"`
		IsPublished bool    `json:"is_published"`
		Expiration  string  `json:"expiration"`
		Price       float64 `json:"price"`
	}

	var r request

	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidProductData.Error(),
		))
		return
	}

	err = ph.svc.UpdateName(id, r.Name)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, web.ErrResponse(
			http.StatusUnprocessableEntity,
			"Unprocessable Entity",
			product.ErrInvalidProductData.Error(),
		))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *productHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"Bad Request",
			product.ErrInvalidId.Error(),
		))
		return
	}

	err = ph.svc.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			err.Error(),
		))
		return
	}

	ctx.Status(http.StatusNoContent)
}
