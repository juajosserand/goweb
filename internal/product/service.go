package product

import "gituhb.com/juajosserand/goweb/internal/domain"

type ProductService interface {
	All() ([]domain.Product, error)
	GetById(int) (domain.Product, error)
	PriceGreaterThan(float64) ([]domain.Product, error)
	Create(string, int, string, bool, string, float64) error
	Update(int, string, int, string, bool, string, float64) error
	UpdateName(int, string) error
	Delete(int) error
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
		return ErrInvalidProductData
	}

	err := p.ToDDMMYYYY()
	if err != nil {
		return ErrInvalidProductData
	}

	err = s.repo.Create(p)
	if err != nil {
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
		return ErrInvalidProductData
	}

	err := p.ToDDMMYYYY()
	if err != nil {
		return ErrInvalidProductData
	}

	err = s.repo.Update(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateName(id int, name string) error {
	if name == "" {
		return ErrInvalidProductData
	}

	err := s.repo.UpdateName(id, name)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
