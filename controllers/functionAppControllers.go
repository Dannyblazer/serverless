package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"serverless/initializers"
	"serverless/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FunctionAppList struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	if filepath.Ext(funcFile.Filename) != ".go" {
		c.JSON(400, gin.H{"error": "Only .go files allowed"})
		return
	}
	// Save file to functions directory
	savePath := filepath.Join("functions", funcName+".go")
	if err := c.SaveUploadedFile(funcFile, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save function file"})
		return
	}
	fmt.Printf("here %v", funcFile.Size)

	// Save function App name and path to DB
	functionApp := models.FunctionApp{
		Name: funcName,
		Path: savePath,
		Size: funcFile.Size,
	}
	if err := initializers.DB.Create(&functionApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error Deploying Function App",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Function deployed successfully", "path": savePath})
}

func ListFunction(c *gin.Context) {
	// Fetch all available function apps
	var functionApps []models.FunctionApp

	if err := initializers.DB.Find(&functionApps).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error Retrieving Function Apps",
		})
		return
	}

	data := make([]FunctionAppList, len(functionApps))
	for i, functionApp := range functionApps {
		data[i] = FunctionAppList{
			ID:        functionApp.ID,
			Name:      functionApp.Name,
			Path:      functionApp.Path,
			CreatedAt: functionApp.CreatedAt,
			UpdatedAt: functionApp.CreatedAt,
		}
	}

	// return list
	fmt.Println("Function apps listed")
	c.JSON(http.StatusOK, gin.H{"message": "Function listed successfully", "data": data})
}

func InvokeFunction(c *gin.Context) {
	// fetch the function app by name/ID
	idStr := c.Param(("id"))
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid FunctionApp ID",
		})
		return
	}

	// Retrieve functionApp from DB and store result in variable
	var functionApp models.FunctionApp
	if err := initializers.DB.Select("id", "name", "path", "size").
		First(&functionApp, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Failed to fetch functionApp",
			})
			return
		}
	}

	// launch the function app
	if _, err := os.Stat(functionApp.Path); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Function App path does not exists",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "run", functionApp.Path)
	output, err := command.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failure during function App invokation",
		})
		return
	}
	fmt.Println(string(output))

	c.JSON(http.StatusOK, gin.H{
		"response": "Function App Invoked Successfully",
	})

}
