// handlers/item_handler.go
package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"examples.com/go-restapi/middleware"
	"examples.com/go-restapi/models"
)

func SetupItemRoutes(router *gin.Engine) {
	api := router.Group("/api/items")
	api.Use(middleware.VerifyToken)
	{
		api.POST("/", CreateItem)
		api.GET("/:id", GetItem)
		api.PUT("/:id", UpdateItem)
		api.DELETE("/:id", DeleteItem)
	}
}

func CreateItem(c *gin.Context) {
	var newItem models.Item
	err := c.BindJSON(&newItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem.ID = uuid.New().String()

	err = saveItem(newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save item"})
		return
	}

	c.JSON(http.StatusOK, newItem)
}

func GetItem(c *gin.Context) {
	itemID := c.Param("id")
	item, err := loadItem(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func UpdateItem(c *gin.Context) {
	itemID := c.Param("id")

	var updatedItem models.Item
	err := c.BindJSON(&updatedItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedItem.ID != itemID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID in URL does not match ID in body"})
		return
	}

	err = saveItem(updatedItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func DeleteItem(c *gin.Context) {
	itemID := c.Param("id")
	err := deleteItem(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func saveItem(item models.Item) error {
	storagePath := viper.GetString("storage_path")
	filePath := filepath.Join(storagePath, "items", item.ID+".json")

	err := writeJSONFile(filePath, item)
	if err != nil {
		return err
	}

	return nil
}

func loadItem(id string) (models.Item, error) {
	storagePath := viper.GetString("storage_path")
	filePath := filepath.Join(storagePath, "items", id+".json")

	var item models.Item
	err := readJSONFile(filePath, &item)
	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func deleteItem(id string) error {
	storagePath := viper.GetString("storage_path")
	filePath := filepath.Join(storagePath, "items", id+".json")

	err := deleteFile(filePath)
	if err != nil {
		return err
	}

	return nil
}

func writeJSONFile(filePath string, data interface{}) error {
	itemJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, itemJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readJSONFile(filePath string, data interface{}) error {
	itemJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemJSON, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
