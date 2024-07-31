// handlers/cache_handler.go
package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"examples.com/go-restapi/models"
	"github.com/spf13/viper"
)

// CacheItem copies an item from storage to cache.
func CacheItem(itemID string) error {
	storagePath := viper.GetString("storage_path")
	cachePath := viper.GetString("cache_path")
	itemFilePath := filepath.Join(storagePath, "items", itemID+".json")
	cacheFilePath := filepath.Join(cachePath, itemID+".json")

	// Create cache directory if it doesn't exist
	err := os.MkdirAll(cachePath, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create cache directory: %v\n", err)
		return err
	}

	// Copy item JSON to cache
	err = copyFile(itemFilePath, cacheFilePath)
	if err != nil {
		log.Printf("Failed to copy item to cache: %v\n", err)
		return err
	}

	return nil
}

// LoadItemFromCache retrieves an item from the cache.
func LoadItemFromCache(itemID string) (models.Item, error) {
	cachePath := viper.GetString("cache_path")
	cacheFilePath := filepath.Join(cachePath, itemID+".json")

	// Read JSON from cache file
	itemJSON, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return models.Item{}, err
	}

	// Unmarshal JSON into item struct
	var item models.Item
	err = json.Unmarshal(itemJSON, &item)
	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}

// DeleteItemFromCache removes an item from the cache.
func DeleteItemFromCache(itemID string) error {
	cachePath := viper.GetString("cache_path")
	cacheFilePath := filepath.Join(cachePath, itemID+".json")

	// Remove file from cache
	err := os.Remove(cacheFilePath)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = sourceFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
