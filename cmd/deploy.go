package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your serverless application",
	Long:  `Deploy your serverless application to the cloud provider.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Deployment logic goes here
		fmt.Println("Deploying your serverless application...")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
