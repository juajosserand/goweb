package product

import "errors"

var (
	errDuplicatedCodeValue = errors.New("duplicated product code value")
	errInvalidId           = errors.New("invalid product id")
	errInvalidPrice        = errors.New("invalid product price")
	errCreation            = errors.New("unable to create product")
	errInvalidProductData  = errors.New("invalid product data")
)
