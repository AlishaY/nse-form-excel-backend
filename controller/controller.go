// controller/controller.go
package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "nse-form-excel-backend/models"
    "gorm.io/gorm"
)

func GetData(c *gin.Context, db *gorm.DB) {
	filepaths := []models.Docs{}

	// Retrieve all files from the database
	if err := db.Find(&filepaths).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving files from database", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files fetched successfully", "data": filepaths})
}