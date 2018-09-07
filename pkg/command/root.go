package command

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "token-tool",
	Short: "Tool for authentication tokens from Keystone",
	Long:  "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()

}
