package repository

import (
	"gituhb.com/juajosserand/goweb/model"
)

func (r *repository) PriceGreaterThan(price float64) ([]model.Product, error) {
	if price < 0 {
		return []model.Product{}, errInvalidPrice
	}

	var products []model.Product
	for _, p := range r.Products {
		if p.Price > price {
			products = append(products, p)
		}
	}

	return products, nil
}
