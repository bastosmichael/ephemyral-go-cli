// build.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds your project",
	Long: `A longer description of your build command that spans multiple lines and likely contains
examples and usage of using your build command. For example:

The build command compiles your code and prepares it for deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building the project...")
		// Here you can add your build logic, such as compiling code, copying resources, etc.
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application, and Local Flags, which will
	// only run when this command is called directly.
	buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
