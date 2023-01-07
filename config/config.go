package config

import "os"

type (
	Config struct {
		HTTP HTTP
		File File
	}

	HTTP struct {
		Port string
	}

	File struct {
		Path string
	}
)

func New() (*Config, error) {
	c := &Config{}

	err := setEnv()
	if err != nil {
		return c, err
	}

	c.File.Path = os.Getenv("PRODUCTS_FILENAME")
	c.HTTP.Port = os.Getenv("HTTP_SERVER_PORT")

	return c, nil
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
