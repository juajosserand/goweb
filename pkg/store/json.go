package store

import (
	"encoding/json"
	"fmt"
	"os"
)

func readJSON(path string, dest any) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return fmt.Errorf("[store.ReadJSON] error: %w", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&dest)
	if err != nil {
		return fmt.Errorf("[store.ReadJSON] error: %w", err)
	}

	return nil
}

func writeJSON(path string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("[store.WriteJSON] error: %w", err)
	}

	err = os.WriteFile(path, jsonData, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}
