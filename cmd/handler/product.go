package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	producti "gituhb.com/juajosserand/goweb/internal/product"
	"gituhb.com/juajosserand/goweb/pkg/storage"
	"gituhb.com/juajosserand/goweb/pkg/web"
)

type product struct {
	svc producti.ProductService
}

func NewProduct(mux *gin.Engine, s producti.ProductService) {
	ph := &product{
		svc: s,
	}

	mux.GET("/ping", ph.Pong)

	productsMux := mux.Group("/products")
	productsMux.GET("/", ph.GetAll)
	productsMux.GET("/:id", ph.GetById)
	productsMux.GET("/search", ph.Search)
	productsMux.GET("/consumer_price", ph.ConsumerPrice)

	productsMux.Use(auth)

	productsMux.POST("/", ph.Create)
	productsMux.PUT("/:id", ph.Update)
	productsMux.PATCH("/:id", ph.PartialUpdate)
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

func auth(ctx *gin.Context) {
	if ctx.GetHeader("token") != os.Getenv("TOKEN") {
		ctx.JSON(http.StatusUnauthorized, web.ErrResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			"invalid token",
		))
		return
	}

	ctx.Next()
}

func (ph *product) Pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, web.Response("pong"))
}

func (ph *product) GetAll(ctx *gin.Context) {
	ps, err := ph.svc.All()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
			http.StatusInternalServerError,
			"internal server error",
			"internal server error",
		))
		return
	}

	ctx.JSON(http.StatusOK, web.Response(ps))
}

func (ph *product) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidId.Error(),
		))
		return
	}

	p, err := ph.svc.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, producti.ErrNotFound):
			ctx.JSON(http.StatusNotFound, web.ErrResponse(
				http.StatusNotFound,
				"not found",
				producti.ErrNotFound.Error(),
			))
		}
		return
	}

	ctx.JSON(http.StatusFound, web.Response(p))
}

func (ph *product) Search(ctx *gin.Context) {
	price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidPrice.Error(),
		))
		return
	}

	ps, err := ph.svc.PriceGreaterThan(price)
	if err != nil {
		switch {
		case errors.Is(err, producti.ErrInvalidPrice):
			ctx.JSON(http.StatusBadRequest, web.ErrResponse(
				http.StatusBadRequest,
				"bad request",
				producti.ErrInvalidPrice.Error(),
			))
		}
		return
	}

	ctx.JSON(http.StatusOK, web.Response(ps))
}

func (ph *product) Create(ctx *gin.Context) {
	var r request

	err := ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidData.Error(),
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
		switch {
		case errors.Is(err, producti.ErrInvalidData):
			ctx.JSON(http.StatusBadRequest, web.ErrResponse(
				http.StatusBadRequest,
				"bad request",
				producti.ErrInvalidData.Error(),
			))
			log.Println(err)
		case errors.Is(err, producti.ErrDuplicatedCodeValue):
			ctx.JSON(http.StatusUnprocessableEntity, web.ErrResponse(
				http.StatusUnprocessableEntity,
				"unprocessable entity",
				producti.ErrDuplicatedCodeValue.Error(),
			))
		case errors.Is(err, storage.ErrWriteFile):
			ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
				http.StatusInternalServerError,
				"internal server error",
				producti.ErrCreation.Error(),
			))
		}
		return
	}

	ctx.Status(http.StatusCreated)
}

func (ph *product) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidId.Error(),
		))
		return
	}

	var r request

	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidData.Error(),
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
		switch {
		case errors.Is(err, producti.ErrInvalidData):
			ctx.JSON(http.StatusBadRequest, web.ErrResponse(
				http.StatusBadRequest,
				"bad request",
				producti.ErrInvalidData.Error(),
			))
		case errors.Is(err, producti.ErrDuplicatedCodeValue):
			ctx.JSON(http.StatusUnprocessableEntity, web.ErrResponse(
				http.StatusUnprocessableEntity,
				"unprocessable entity",
				producti.ErrDuplicatedCodeValue.Error(),
			))
		case errors.Is(err, storage.ErrWriteFile):
			ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
				http.StatusInternalServerError,
				"internal server error",
				producti.ErrCreation.Error(),
			))
		case errors.Is(err, producti.ErrNotFound):
			ctx.JSON(http.StatusNotFound, web.ErrResponse(
				http.StatusNotFound,
				"not found",
				producti.ErrNotFound.Error(),
			))
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *product) PartialUpdate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidId.Error(),
		))
		return
	}

	p, err := ph.svc.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, web.ErrResponse(
			http.StatusNotFound,
			"not found",
			producti.ErrNotFound.Error(),
		))
		return
	}

	err = json.NewDecoder(ctx.Request.Body).Decode(&p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidData.Error(),
		))
		return
	}

	err = ph.svc.Update(
		p.Id,
		p.Name,
		p.Quantity,
		p.CodeValue,
		p.IsPublished,
		p.Expiration,
		p.Price,
	)
	if err != nil {
		switch {
		case errors.Is(err, producti.ErrInvalidData):
			ctx.JSON(http.StatusBadRequest, web.ErrResponse(
				http.StatusBadRequest,
				"bad request",
				producti.ErrInvalidData.Error(),
			))
		case errors.Is(err, producti.ErrDuplicatedCodeValue):
			ctx.JSON(http.StatusUnprocessableEntity, web.ErrResponse(
				http.StatusUnprocessableEntity,
				"unprocessable entity",
				producti.ErrDuplicatedCodeValue.Error(),
			))
		case errors.Is(err, storage.ErrWriteFile):
			ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
				http.StatusInternalServerError,
				"internal server error",
				producti.ErrCreation.Error(),
			))
		case errors.Is(err, producti.ErrNotFound):
			ctx.JSON(http.StatusNotFound, web.ErrResponse(
				http.StatusNotFound,
				"not found",
				producti.ErrNotFound.Error(),
			))
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *product) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidId.Error(),
		))
		return
	}

	err = ph.svc.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrWriteFile):
			ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
				http.StatusInternalServerError,
				"internal server error",
				producti.ErrCreation.Error(),
			))
		case errors.Is(err, producti.ErrNotFound):
			ctx.JSON(http.StatusNotFound, web.ErrResponse(
				http.StatusNotFound,
				"not found",
				producti.ErrNotFound.Error(),
			))
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ph *product) ConsumerPrice(ctx *gin.Context) {
	// compile regex
	r, err := regexp.Compile(`\[\d+(?:,\d+)*\]`)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
			http.StatusInternalServerError,
			"internal server error",
			"unable to validate list of product ids",
		))
		return
	}

	// validate list string
	listStr := ctx.Query("list")

	if !r.MatchString(listStr) {
		ctx.JSON(http.StatusBadRequest, web.ErrResponse(
			http.StatusBadRequest,
			"bad request",
			producti.ErrInvalidConsumerPriceList.Error(),
		))
		return
	}

	// parse ids
	listStr = listStr[1 : len(listStr)-1]
	split := strings.Split(listStr, ",")

	// convert ids and count products
	productQuantities := make(map[int]int)

	for _, s := range split {
		id, err := strconv.Atoi(s)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, web.ErrResponse(
				http.StatusInternalServerError,
				"internal server error",
				producti.ErrInvalidConsumerPriceList.Error(),
			))
		}

		productQuantities[id]++
	}

	// compute total
	total, products, err := ph.svc.CustomerPrice(productQuantities)
	if err != nil {
		switch {
		case errors.Is(err, producti.ErrNotFound):
			ctx.JSON(http.StatusNotFound, web.ErrResponse(
				http.StatusNotFound,
				"not found",
				producti.ErrNotFound.Error(),
			))
		case errors.Is(err, producti.ErrNoStock):
			ctx.JSON(http.StatusBadRequest, web.ErrResponse(
				http.StatusBadRequest,
				"bad request",
				producti.ErrNoStock.Error(),
			))
		case errors.Is(err, producti.ErrNotPublished):
			ctx.JSON(http.StatusBadRequest, web.ErrResponse(
				http.StatusBadRequest,
				"bad request",
				producti.ErrNotPublished.Error(),
			))
		}
		return
	}

	ctx.JSON(http.StatusOK, web.Response(gin.H{
		"products":    products,
		"total_price": total,
	}))
}
