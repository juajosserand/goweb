package product

import (
	"encoding/json"
	"fmt"
	"os"

	"gituhb.com/juajosserand/goweb/internal/domain"
)

type ProductRepository interface {
	All() ([]domain.Product, error)
	GetById(int) (domain.Product, error)
	PriceGreaterThan(float64) ([]domain.Product, error)
	Create(domain.Product) error
	Update(domain.Product) error
	UpdateName(int, string) error
	Delete(int) error
}

type repository struct {
	Products []domain.Product `json:"products"`
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

func (r *repository) All() ([]domain.Product, error) {
	return r.Products, nil
}

func (r *repository) GetById(id int) (domain.Product, error) {
	if id < 1 {
		return domain.Product{}, ErrInvalidId
	}

	for _, p := range r.Products {
		if p.Id == id {
			return p, nil
		}
	}

	return domain.Product{}, nil
}

func (r *repository) PriceGreaterThan(price float64) ([]domain.Product, error) {
	if price < 0 {
		return []domain.Product{}, ErrInvalidPrice
	}

	var products []domain.Product
	for _, p := range r.Products {
		if p.Price > price {
			products = append(products, p)
		}
	}

	return products, nil
}

func (r *repository) Create(p domain.Product) error {
	for _, product := range r.Products {
		if product.CodeValue == p.CodeValue {
			return ErrDuplicatedCodeValue
		}
	}

	r.lastId++
	p.Id = r.lastId
	r.Products = append(r.Products, p)
	return nil
}

func (r *repository) Update(p domain.Product) error {
	// find product
	for i, product := range r.Products {
		if product.Id == p.Id {
			// check code value
			for _, pCheck := range r.Products {
				if pCheck.CodeValue == p.CodeValue && p.CodeValue != product.CodeValue {
					return ErrDuplicatedCodeValue
				}
			}

			r.Products[i] = p

			return nil
		}
	}

	return ErrUnexistingProduct
}

func (r *repository) UpdateName(id int, name string) error {
	for i, product := range r.Products {
		if product.Id == id {
			r.Products[i].Name = name
			return nil
		}
	}

	return ErrUnexistingProduct
}

func (r *repository) Delete(id int) error {
	for i, product := range r.Products {
		if product.Id == id {
			r.Products = append(r.Products[:i], r.Products[i+1:]...)
			return nil
		}
	}

	return ErrUnexistingProduct
}
