package domain

import (
	"fmt"
	"time"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,gte=1"`
	CodeValue   string  `json:"code_value" validate:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" validate:"required"`
	Price       float64 `json:"price" validate:"required,gte=0"`
}

func (p *Product) IsExpirationValid() bool {
	expDate, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {
		return false
	}

	return !expDate.Before(time.Now())
}

func (p *Product) ToDDMMYYYY() error {
	expDate, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {
		return err
	}

	p.Expiration = ""

	if expDate.Day() < 10 {
		p.Expiration += fmt.Sprintf("0%d/", expDate.Day())
	} else {
		p.Expiration += fmt.Sprintf("%d/", expDate.Day())
	}

	if expDate.Month() < 10 {
		p.Expiration += fmt.Sprintf("0%d/", expDate.Month())
	} else {
		p.Expiration += fmt.Sprintf("%d/", expDate.Month())
	}

	p.Expiration += fmt.Sprintf("%d", expDate.Year())

	return nil
}
