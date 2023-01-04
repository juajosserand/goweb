package service

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gituhb.com/juajosserand/goweb/repository"
)

var (
	errInvalidId    = errors.New("invalid product id")
	errInvalidPrice = errors.New("invalid product price")
)

type server struct {
	repo repository.ProductsRepository
	mux  *gin.Engine
	port string
}

func New(r repository.ProductsRepository, p string) *server {
	return &server{
		repo: r,
		mux:  gin.Default(),
		port: fmt.Sprintf(":%s", p),
	}
}

func (s *server) Run() error {
	s.mux.GET("/ping", s.pong)
	s.mux.GET("/products", s.getAll)
	s.mux.GET("/products/:id", s.getById)
	s.mux.GET("/products/search", s.search)

	if err := s.mux.Run(s.port); err != nil {
		return err
	}

	return nil
}
