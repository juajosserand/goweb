package repository

import "gituhb.com/juajosserand/goweb/model"

func (r *repository) Create(p model.Product) {
	r.lastId++
	p.Id = r.lastId
	r.Products = append(r.Products, p)
}
