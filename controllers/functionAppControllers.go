package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DeployFunction(c *gin.Context) {
	// Parse multipart form
	funcFile, err := c.FormFile("function")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function file is required"})
		return
	}
	funcName := c.PostForm("name")
	if funcName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function name is required"})
		return
	}
	// Save file to functions directory
	savePath := filepath.Join("functions", funcName+".go")
	if err := c.SaveUploadedFile(funcFile, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save function file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Function deployed successfully", "path": savePath})
}
