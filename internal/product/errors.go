package product

import "errors"

var (
	ErrInvalidProductData  = errors.New("invalid product data")
	ErrInvalidId           = errors.New("invalid product id")
	ErrInvalidPrice        = errors.New("invalid product price")
	ErrInvalidToken        = errors.New("invalid token")
	ErrCreation            = errors.New("unable to create product")
	ErrDeletion            = errors.New("unable to delete product")
	ErrUnexistingProduct   = errors.New("unable to find product")
	ErrDuplicatedCodeValue = errors.New("duplicated product code value")
)
