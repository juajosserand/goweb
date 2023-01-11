package storage

import (
	"errors"
	"fmt"
	"path/filepath"

	"gituhb.com/juajosserand/goweb/internal/domain"
)

var (
	ErrReadFile  = errors.New("unable to read file")
	ErrWriteFile = errors.New("unable to write file")
)

func ReadFile(path string, dest any) error {
	ext := filepath.Ext(path)

	switch {
	case ext == ".json":
		err := readJSON(path, dest)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrReadFile, err)
		}

		return nil
	case ext == ".csv":
		products, ok := dest.(*[]domain.Product)
		if !ok {
			return fmt.Errorf("%w: [storage.ReadFile] invalid csv receiver type", ErrReadFile)
		}

		err := readCSV(path, products)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrReadFile, err)
		}

		return nil
	default:
		return fmt.Errorf("%w: [storage.ReadFile] file format %s not supported", ErrReadFile, ext)
	}
}

func WriteFile(path string, data any) error {
	ext := filepath.Ext(path)

	switch {
	case ext == ".json":
		err := writeJSON(path, data)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrWriteFile, err)
		}

		return nil
	case ext == ".csv":
		products, ok := data.(*[]domain.Product)
		if !ok {
			return fmt.Errorf("%w: [storage.WriteFile] invalid csv data type", ErrWriteFile)
		}

		err := writeCSV(path, products)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrWriteFile, err)
		}
		return nil
	default:
		return fmt.Errorf("[storage.WriteFile] file format %s not supported", ext)
	}
}
