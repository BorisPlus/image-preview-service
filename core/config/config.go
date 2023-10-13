package config

import (
	"encoding/json"
	"os"
)

func LoadFromJsonFile(filename string, configStruct any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &configStruct)
}
