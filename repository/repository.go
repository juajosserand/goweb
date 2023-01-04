package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gituhb.com/juajosserand/goweb/model"
)

var (
	errInvalidId    = errors.New("invalid product id")
	errInvalidPrice = errors.New("invalid product price")
)

type ProductRepository interface {
	All() []model.Product
	GetById(int) (model.Product, error)
	PriceGreaterThan(float64) ([]model.Product, error)
}

type repository struct {
	filename string
	Products []model.Product `json:"products"`
}

func New(fn string) (ProductRepository, error) {
	r := &repository{
		filename: fn,
	}

	err := r.readFromFile()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *repository) readFromFile() error {
	f, err := os.OpenFile(r.filename, os.O_RDONLY, 0444)
	if err != nil {
		return fmt.Errorf("[repository.readFromFile] error: %w", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&r.Products)
	if err != nil {
		return fmt.Errorf("[repository.readFromFile] error: %w", err)
	}

	return nil
}
