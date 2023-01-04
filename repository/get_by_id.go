package repository

import (
	"gituhb.com/juajosserand/goweb/model"
)

func (r *repository) GetById(id int) (model.Product, error) {
	if id < 1 {
		return model.Product{}, errInvalidId
	}

	for _, p := range r.Products {
		if p.Id == id {
			return p, nil
		}
	}

	return model.Product{}, nil
}
