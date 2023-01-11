package product

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidData = errors.New("invalid product data")
	ErrCreation    = errors.New("unable to create product")
	ErrDeletion    = errors.New("unable to delete product")
	ErrNotFound    = errors.New("unable to find product")

	ErrInvalidId                = errors.New("invalid product id")
	ErrInvalidPrice             = errors.New("invalid product price")
	ErrDuplicatedCodeValue      = errors.New("duplicated product code value")
	ErrInvalidConsumerPriceList = errors.New("invalid product list")

	ErrInvalidToken = errors.New("invalid token")
)

func ErrNoStock(name string) error {
	return fmt.Errorf("no enough stock for product %s", name)
}

func ErrNotPublished(name string) error {
	return fmt.Errorf("product %s is not published", name)
}
