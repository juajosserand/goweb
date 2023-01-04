package main

import (
	"log"
	"os"

	"gituhb.com/juajosserand/goweb/repository"
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

	_, err = repository.New(os.Getenv("PRODUCTS_FILENAME"))
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
