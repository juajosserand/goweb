package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gituhb.com/juajosserand/goweb/repository"
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

	if err := s.mux.Run(s.port); err != nil {
		return err
	}

	return nil
}
