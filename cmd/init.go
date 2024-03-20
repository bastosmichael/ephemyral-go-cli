// init.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new ephemyral task",
	Long:  `This command initializes a new ephemyral task.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Ephemyral task initialized")
	},
}
