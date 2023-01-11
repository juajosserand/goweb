package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gituhb.com/juajosserand/goweb/cmd/handler"
	"gituhb.com/juajosserand/goweb/internal/product"
	"gituhb.com/juajosserand/goweb/pkg/httpserver"
)

func main() {
	// load env
	err := godotenv.Load()
	if err != nil {
		log.Println(fmt.Errorf("error: %w", err))
	}

	// repository
	repo, err := product.NewRepository()
	if err != nil {
		log.Println(fmt.Errorf("error: %w", err))
	}

	// service
	svc := product.NewService(repo)

	// http server
	mux := gin.Default()
	handler.NewProduct(mux, svc)
	server := httpserver.New(mux, httpserver.Port(os.Getenv("HTTP_SERVER_PORT")))

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
