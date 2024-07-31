package data

import (
	"encoding/json"
	"os"
)

// Example data structure
type Data struct {
	Message string `json:"message"`
}

// Function to read data from JSON file
func ReadDataFromFile(filename string) (Data, error) {
	var data Data

	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()

	// Decode JSON from file into data struct
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

// Function to write data to JSON file
func WriteDataToFile(data Data, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode data struct to JSON and write to file
	if err := json.NewEncoder(file).Encode(data); err != nil {
		return err
	}

	return nil
}
