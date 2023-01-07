package product

import (
	"fmt"
	"time"
)

type ProductService interface {
	All() []Product
	GetById(int) (Product, error)
	PriceGreaterThan(float64) ([]Product, error)
	Create(string, int, string, bool, string, float64) error
}

type service struct {
	repo ProductRepository
}

func NewService(r ProductRepository) ProductService {
	return &service{
		repo: r,
	}
}

func (s *service) All() []Product {
	return s.repo.All()
}

func (s *service) GetById(id int) (Product, error) {
	return s.repo.GetById(id)
}

func (s *service) PriceGreaterThan(p float64) ([]Product, error) {
	return s.repo.PriceGreaterThan(p)
}

func (s *service) Create(name string, quantity int, codeValue string, isPublished bool, expiration string, price float64) error {
	p := Product{
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}

	// expiration validation
	expDate, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {
		return err
	}

	if expDate.Before(time.Now()) {
		return errInvalidProductData
	}

	var dateStr string

	if expDate.Day() < 10 {
		dateStr += fmt.Sprintf("0%d/", expDate.Day())
	} else {
		dateStr += fmt.Sprintf("%d/", expDate.Day())
	}

	if expDate.Month() < 10 {
		dateStr += fmt.Sprintf("0%d/", expDate.Month())
	} else {
		dateStr += fmt.Sprintf("%d/", expDate.Month())
	}

	dateStr += fmt.Sprintf("%d", expDate.Year())

	p.Expiration = dateStr

	// create product
	err = s.repo.Create(p)
	if err != nil {
		return err
	}

	return nil
}
