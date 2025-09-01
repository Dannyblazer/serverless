package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the version of MyCLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MyCLI v1.0.0")
	},
}
