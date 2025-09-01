package main

import (
	"fmt"
	"os"
	"path/filepath"
	"serverless/controllers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func main() {
	// API Initiailization
	router := gin.Default()

	// API Routes
	router.POST("deploy/", controllers.DeployFunction)

	var rootCmd = &cobra.Command{Use: "platform"}
	var funcName, funcPath string

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a new function",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate file exists
			if _, err := os.Stat(funcPath); os.IsNotExist(err) {
				fmt.Printf("Error: File %s does not exist\n", funcPath)
				return
			}
			// Copy file to functions directory
			destPath := filepath.Join("functions", funcName+".go")
			input, err := os.ReadFile(funcPath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}
			err = os.WriteFile(destPath, input, 0644)
			if err != nil {
				fmt.Printf("Error saving function: %v\n", err)
				return
			}
			fmt.Printf("Deployed function %s to %s\n", funcName, destPath)
		},
	}
	deployCmd.Flags().StringVar(&funcName, "name", "", "Function name")
	deployCmd.Flags().StringVar(&funcPath, "path", "", "Path to function file")
	deployCmd.MarkFlagRequired("name")
	deployCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(deployCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//cmd.Execute()
}
