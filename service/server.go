package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gituhb.com/juajosserand/goweb/repository"
)

var (
	errInvalidId             = errors.New("invalid product id")
	errInvalidPrice          = errors.New("invalid product price")
	errDuplicatedCodeValue   = errors.New("duplicated product code value")
	errInvalidExpirationDate = errors.New("expiration date cannot be equal or less than today")
)

type server struct {
	repo repository.ProductRepository
	mux  *gin.Engine
	port string
}

func New(r repository.ProductRepository, p string) *server {
	return &server{
		repo: r,
		mux:  gin.Default(),
		port: fmt.Sprintf(":%s", p),
	}
}

func (s *server) Run() error {
	s.mux.GET("/ping", s.pong)

	productsMux := s.mux.Group("/products")
	productsMux.GET("/", s.getAll)
	productsMux.GET("/:id", s.getById)
	productsMux.GET("/search", s.search)
	productsMux.POST("/", s.create)

	if err := s.mux.Run(s.port); err != nil {
		return err
	}

	return nil
}
