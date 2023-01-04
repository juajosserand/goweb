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

type ProductsRepository interface {
	All() []model.Product
	GetById(int) (model.Product, error)
	PriceGreaterThan(float64) ([]model.Product, error)
}

type repository struct {
	filename string
	Products []model.Product `json:"products"`
}

func New(fn string) (ProductsRepository, error) {
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
	data, err := os.ReadFile(r.filename)
	if err != nil {
		return fmt.Errorf("[repository.readFromFile] error: %w", err)
	}

	err = json.Unmarshal(data, &r.Products)
	if err != nil {
		return fmt.Errorf("[repository.readFromFile] error: %w", err)
	}

	return nil
}
