package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gituhb.com/juajosserand/goweb/cmd/handlers"
)

type server struct {
	handlers handlers.ProductHandlers
	mux      *gin.Engine
	port     string
}

func New(ph handlers.ProductHandlers, p string) *server {
	return &server{
		handlers: ph,
		mux:      gin.Default(),
		port:     fmt.Sprintf(":%s", p),
	}
}

func (s *server) Run() error {
	s.mux.GET("/ping", s.handlers.Pong)

	productsMux := s.mux.Group("/products")
	productsMux.GET("/", s.handlers.GetAll)
	productsMux.GET("/:id", s.handlers.GetById)
	productsMux.GET("/search", s.handlers.Search)
	productsMux.POST("/", s.handlers.Create)

	if err := s.mux.Run(s.port); err != nil {
		return err
	}

	return nil
}
