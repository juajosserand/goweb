package product

import (
	"fmt"
	"time"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (p *Product) IsExpirationValid() bool {
	expDate, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {
		return false
	}

	if expDate.Before(time.Now()) {
		return false
	}

	return true
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
