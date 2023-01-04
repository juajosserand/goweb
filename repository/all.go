package repository

import "gituhb.com/juajosserand/goweb/model"

func (r *repository) All() []model.Product {
	return r.Products
}
