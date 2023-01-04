package main

import (
	"log"
	"os"

	"gituhb.com/juajosserand/goweb/repository"
	"gituhb.com/juajosserand/goweb/service"
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

	r, err := repository.New(os.Getenv("PRODUCTS_FILENAME"))
	if err != nil {
		panic(err)
	}

	s := service.New(r, os.Getenv("HTTP_SERVER_PORT"))
	err = s.Run()
	if err != nil {
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
