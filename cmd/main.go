package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"gituhb.com/juajosserand/goweb/config"
	"gituhb.com/juajosserand/goweb/internal/product"
	"gituhb.com/juajosserand/goweb/pkg/httpserver"
)

func main() {
	// reading config
	config, err := config.New()
	if err != nil {
		log.Println(fmt.Errorf("error: %w", err))
	}

	// repository
	repo, err := product.NewRepository(config.File.Path)
	if err != nil {
		log.Println(fmt.Errorf("error: %w", err))
	}

	// service
	svc := product.NewService(repo)

	// http server
	mux := gin.Default()
	product.NewHandler(mux, svc)

	server := httpserver.New(mux, httpserver.Port(config.HTTP.Port))

	// signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("signal:", s.String())
	case err := <-server.Notify():
		log.Println(fmt.Errorf("error: %w", err))
	}

	err = server.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("error: %w", err))
	}
}
