package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	// Create the main "platform" command
	var rootCmd = &cobra.Command{Use: "platform"}

	// Variables to store flag values
	var funcName, funcPath string

	// Create the "deploy" subcommand
	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a new function",
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the file exists
			if _, err := os.Stat(funcPath); os.IsNotExist(err) {
				fmt.Printf("Error: File %s does not exist\n", funcPath)
				return
			}
			// Create functions directory if it doesn't exist
			os.MkdirAll("functions", 0755)
			// Copy the file to functions directory
			destPath := filepath.Join("functions", funcName+".go")
			input, err := os.ReadFile(funcPath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}
			err = os.WriteFile(destPath, input, 0644)
			if err != nil {
				fmt.Printf("Error saving file: %v\n", err)
				return
			}
			fmt.Printf("Deployed function %s to %s\n", funcName, destPath)
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available functions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello People")
			// List function Logic goes here
		},
	}

	var invokeCmd = &cobra.Command{
		Use:   "Invoke",
		Short: "Invoke an app function using its name",
		Run: func(cmd *cobra.Command, args []string) {
			funcPath := filepath.Join("functions", funcName+".go")
			if _, err := os.Stat(funcPath); os.IsNotExist(err) {
				fmt.Printf("Error: File %s does not exist\n", funcPath)
				return
			}
			command := exec.Command("go", "run", funcPath)
			output, err := command.CombinedOutput()
			if err != nil {
				fmt.Printf("Error running file: %v\n", err)
				return
			}

			fmt.Println(string(output))
		},
	}

	// Add flags for function name and path
	deployCmd.Flags().StringVar(&funcName, "name", "", "Function name")
	deployCmd.Flags().StringVar(&funcPath, "path", "", "Path to function file")
	deployCmd.MarkFlagRequired("name")
	deployCmd.MarkFlagRequired("path")

	// Add funcName to Invoke command
	invokeCmd.Flags().StringVar(&funcName, "name", "", "Function name")
	invokeCmd.MarkFlagRequired("name")

	// Add deploy command to root
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(invokeCmd)

	// Run the CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
