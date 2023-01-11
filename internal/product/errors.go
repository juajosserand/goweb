package product

import (
	"errors"
)

var (
	ErrInvalidData = errors.New("invalid product data")
	ErrCreation    = errors.New("unable to create product")
	ErrDeletion    = errors.New("unable to delete product")
	ErrNotFound    = errors.New("unable to find product")

	ErrInvalidId                = errors.New("invalid product id")
	ErrInvalidPrice             = errors.New("invalid product price")
	ErrDuplicatedCodeValue      = errors.New("duplicated product code value")
	ErrInvalidConsumerPriceList = errors.New("invalid list of product ids")
	ErrNoStock                  = errors.New("no enough stock for product")
	ErrNotPublished             = errors.New("not published product")
)
