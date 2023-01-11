package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"gituhb.com/juajosserand/goweb/internal/domain"
)

func readCSV(path string, dest *[]domain.Product) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("[storage.readCSV] error: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("[storage.readCSV] error: %w", err)
	}

	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return fmt.Errorf("[storage.readCSV] error: %w", err)
		}

		quantity, err := strconv.Atoi(record[2])
		if err != nil {
			return fmt.Errorf("[storage.readCSV] error: %w", err)
		}

		isPublished, err := strconv.ParseBool(record[4])
		if err != nil {
			return fmt.Errorf("[storage.readCSV] error: %w", err)
		}

		price, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return fmt.Errorf("[storage.readCSV] error: %w", err)
		}

		*dest = append(*dest, domain.Product{
			Id:          id,
			Name:        record[1],
			Quantity:    quantity,
			CodeValue:   record[3],
			IsPublished: isPublished,
			Expiration:  record[5],
			Price:       price,
		})
	}

	return nil
}

func writeCSV(path string, data *[]domain.Product) error {
	if _, err := os.Stat(path); err == nil {
		// file exist
		os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("[storage.writeCSV] error: %w", err)
	}
	defer f.Close()

	var records [][]string
	for _, p := range *data {
		records = append(records, []string{
			strconv.Itoa(p.Id),
			p.Name,
			strconv.Itoa(p.Quantity),
			p.CodeValue,
			strconv.FormatBool(p.IsPublished),
			p.Expiration,
			strconv.FormatFloat(p.Price, 'e', 2, 64),
		})
	}

	writer := csv.NewWriter(f)
	err = writer.WriteAll(records)
	if err != nil {
		return fmt.Errorf("[storage.writeCSV] error: %w", err)
	}

	return nil
}
