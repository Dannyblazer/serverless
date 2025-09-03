package main

import (
	"serverless/initializers"
	"serverless/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.FunctionApp{})
}
