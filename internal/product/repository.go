package product

import (
	"encoding/json"
	"fmt"
	"os"
)

type ProductRepository interface {
	All() []Product
	GetById(int) (Product, error)
	PriceGreaterThan(float64) ([]Product, error)
	Create(Product) error
}

type repository struct {
	Products []Product `json:"products"`
	filename string
	lastId   int
}

func NewRepository(fn string) (ProductRepository, error) {
	r := &repository{
		filename: fn,
	}

	err := r.readFromFile()
	if err != nil {
		return nil, err
	}

	r.lastId = r.Products[len(r.Products)-1].Id

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

func (r *repository) All() []Product {
	return r.Products
}

func (r *repository) GetById(id int) (Product, error) {
	if id < 1 {
		return Product{}, errInvalidId
	}

	for _, p := range r.Products {
		if p.Id == id {
			return p, nil
		}
	}

	return Product{}, nil
}

func (r *repository) PriceGreaterThan(price float64) ([]Product, error) {
	if price < 0 {
		return []Product{}, errInvalidPrice
	}

	var products []Product
	for _, p := range r.Products {
		if p.Price > price {
			products = append(products, p)
		}
	}

	return products, nil
}

func (r *repository) Create(p Product) error {
	for _, product := range r.Products {
		if product.CodeValue == p.CodeValue {
			return errDuplicatedCodeValue
		}
	}

	r.lastId++
	p.Id = r.lastId
	r.Products = append(r.Products, p)
	return nil
}
