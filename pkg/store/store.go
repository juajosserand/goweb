package store

import (
	"fmt"
	"path/filepath"

	"gituhb.com/juajosserand/goweb/internal/domain"
)

func ReadFile(path string, dest any) error {
	ext := filepath.Ext(path)

	switch {
	case ext == ".json":
		return readJSON(path, dest)
	case ext == ".csv":
		products, ok := dest.(*[]domain.Product)
		if !ok {
			return fmt.Errorf("[store.ReadFile] invalid receiver for csv format")
		}

		return readCSV(path, products)
	default:
		return fmt.Errorf("[store.ReadFile] file format %s not supported", ext)
	}
}

func WriteFile(path string, data any) error {
	ext := filepath.Ext(path)

	switch {
	case ext == ".json":
		return writeJSON(path, data)
	case ext == ".csv":
		products, ok := data.(*[]domain.Product)
		if !ok {
			return fmt.Errorf("[store.WriteFile] invalid data for csv format")
		}

		return writeCSV(path, products)
	default:
		return fmt.Errorf("[store.WriteFile] file format %s not supported", ext)
	}
}
