package datagen

import (
	"encoding/json"
	"os"
)

// Options defines the data generator options
type Options struct {
	Users uint64 `json:"users"`
	Posts uint64 `json:"posts"`
}

// OptionsFromFile loads the options from a JSON file
func OptionsFromFile(filePath string) (options Options, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&options)

	return
}
