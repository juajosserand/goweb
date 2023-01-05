package main

import (
	"log"
	"os"

	"gituhb.com/juajosserand/goweb/cmd/handlers"
	"gituhb.com/juajosserand/goweb/cmd/server"
	"gituhb.com/juajosserand/goweb/internal/product"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	err := setEnv()
	if err != nil {
		panic(err)
	}

	repo, err := product.NewRepository(os.Getenv("PRODUCTS_FILENAME"))
	if err != nil {
		panic(err)
	}

	svc, err := product.NewService(repo)
	if err != nil {
		panic(err)
	}

	ph, err := handlers.NewProductHandlers(svc)
	if err != nil {
		panic(err)
	}

	s := server.New(ph, os.Getenv("HTTP_SERVER_PORT"))
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func setEnv() (err error) {
	err = os.Setenv("PRODUCTS_FILENAME", "./products.json")
	if err != nil {
		return
	}

	err = os.Setenv("HTTP_SERVER_PORT", "8080")
	if err != nil {
		return
	}

	return
}
