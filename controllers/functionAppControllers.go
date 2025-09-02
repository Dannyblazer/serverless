package controllers

import (
	"fmt"
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

func ListFunction(c *gin.Context) {
	// Fetch all available function apps

	// return list
	fmt.Println("Function apps listed")
	c.JSON(http.StatusOK, gin.H{"message": "Function listed successfully"})
}

func InvokeFunction(c *gin.Context) {
	// fetch the function app by name/ID
	// launch the function app
	var funcPath string
	fmt.Println(funcPath)

}
