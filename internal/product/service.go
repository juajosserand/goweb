package product

import (
	"fmt"

	"gituhb.com/juajosserand/goweb/internal/domain"
)

type ProductService interface {
	All() ([]domain.Product, error)
	GetById(int) (domain.Product, error)
	PriceGreaterThan(float64) ([]domain.Product, error)
	Create(string, int, string, bool, string, float64) error
	Update(int, string, int, string, bool, string, float64) error
	PartialUpdate(int, string, int, string, bool, string, float64) error
	Delete(int) error
	CustomerPrice(map[int]int) (float64, []domain.Product, error)
}

type service struct {
	repo ProductRepository
}

func NewService(r ProductRepository) ProductService {
	return &service{
		repo: r,
	}
}

func (s *service) All() ([]domain.Product, error) {
	products, err := s.repo.All()
	if err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (s *service) GetById(id int) (domain.Product, error) {
	return s.repo.GetById(id)
}

func (s *service) PriceGreaterThan(p float64) ([]domain.Product, error) {
	if p < 0 {
		return []domain.Product{}, ErrInvalidPrice
	}

	return s.repo.PriceGreaterThan(p)
}

func (s *service) Create(name string, quantity int, codeValue string, isPublished bool, expiration string, price float64) error {
	p := domain.Product{
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}

	if !p.IsExpirationValid() {
		return ErrInvalidData
	}

	if err := p.ToDDMMYYYY(); err != nil {
		return ErrInvalidData
	}

	if err := s.repo.Create(p); err != nil {
		return err
	}

	return nil
}

func (s *service) Update(id int, name string, quantity int, codeValue string, isPublished bool, expiration string, price float64) error {
	p := domain.Product{
		Id:          id,
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}

	if !p.IsExpirationValid() {
		return ErrInvalidData
	}

	err := p.ToDDMMYYYY()
	if err != nil {
		return ErrInvalidData
	}

	err = s.repo.Update(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) PartialUpdate(id int, name string, quantity int, codeValue string, isPublished bool, expiration string, price float64) error {
	p := domain.Product{
		Id:          id,
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}

	err := s.repo.Update(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *service) CustomerPrice(quantities map[int]int) (total float64, products []domain.Product, err error) {
	var numProducts int

	for id, q := range quantities {
		p, err := s.repo.GetById(id)
		if err != nil {
			return total, products, err
		}

		if q > p.Quantity {
			err = fmt.Errorf("%w: %s", ErrNoStock, p.Name)
			return total, products, err
		}

		if !p.IsPublished {
			err = fmt.Errorf("%w: %s", ErrNoStock, p.Name)
			return total, products, err
		}

		products = append(products, p)
		numProducts += q
		total += p.Price * float64(q)
	}

	switch {
	case numProducts < 10:
		total *= 1.21
	case numProducts >= 10 && numProducts < 20:
		total *= 1.17
	case numProducts >= 20:
		total *= 1.15
	}

	return
}
