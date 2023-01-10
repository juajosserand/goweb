package product

type ProductService interface {
	All() []Product
	GetById(int) (Product, error)
	PriceGreaterThan(float64) ([]Product, error)
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

	if !p.IsExpirationValid() {
		return errInvalidProductData
	}

	err := p.ToDDMMYYYY()
	if err != nil {
		return errInvalidProductData
	}

	err = s.repo.Create(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Update(id int, name string, quantity int, codeValue string, isPublished bool, expiration string, price float64) error {
	p := Product{
		Id:          id,
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}

	if !p.IsExpirationValid() {
		return errInvalidProductData
	}

	err := p.ToDDMMYYYY()
	if err != nil {
		return errInvalidProductData
	}

	err = s.repo.Update(p)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateName(id int, name string) error {
	if name == "" {
		return errInvalidProductData
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
