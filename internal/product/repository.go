package product

import (
	"os"

	"gituhb.com/juajosserand/goweb/internal/domain"
	"gituhb.com/juajosserand/goweb/pkg/storage"
)

type ProductRepository interface {
	All() ([]domain.Product, error)
	GetById(int) (domain.Product, error)
	PriceGreaterThan(float64) ([]domain.Product, error)
	Create(domain.Product) error
	Update(domain.Product) error
	Delete(int) error
}

type repository struct {
	Products []domain.Product `json:"products"`
	lastId   int
}

func NewRepository() (ProductRepository, error) {
	r := &repository{}

	err := storage.ReadFile(os.Getenv("PRODUCTS_FILENAME"), &r.Products)
	if err != nil {
		return r, err
	}

	if len(r.Products) > 0 {
		r.lastId = r.Products[len(r.Products)-1].Id
	}

	return r, nil
}

func (r *repository) All() ([]domain.Product, error) {
	return r.Products, nil
}

func (r *repository) GetById(id int) (domain.Product, error) {
	for _, p := range r.Products {
		if p.Id == id {
			return p, nil
		}
	}

	return domain.Product{}, ErrNotFound
}

func (r *repository) PriceGreaterThan(price float64) (products []domain.Product, err error) {
	for _, p := range r.Products {
		if p.Price > price {
			products = append(products, p)
		}
	}

	return
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

	err := storage.WriteFile(os.Getenv("PRODUCTS_FILENAME"), &r.Products)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(p domain.Product) error {
	for i, product := range r.Products {
		if product.Id == p.Id {
			// check code value
			for _, pCheck := range r.Products {
				if pCheck.CodeValue == p.CodeValue && p.CodeValue != product.CodeValue {
					return ErrDuplicatedCodeValue
				}
			}

			r.Products[i] = p

			err := storage.WriteFile(os.Getenv("PRODUCTS_FILENAME"), &r.Products)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrNotFound
}

func (r *repository) Delete(id int) error {
	for i, product := range r.Products {
		if product.Id == id {
			r.Products = append(r.Products[:i], r.Products[i+1:]...)

			err := storage.WriteFile(os.Getenv("PRODUCTS_FILENAME"), &r.Products)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrNotFound
}
