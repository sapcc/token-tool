package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VERSION = "dev"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version",
	Long:  "Prints the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}
