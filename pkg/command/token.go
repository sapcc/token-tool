package command

import (
	"github.com/sapcc/token-tool/pkg/legacy"
	"github.com/spf13/cobra"
	"log"
)

func init() {

	var tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "Retrieves token from Keystone and print",
		Long:  ``,
		Run:   Token,
	}

	var keystoneEndpoint, user, password, userDomainName, project, projectDomainName, format string

	//TODO bind to environment variables via viper
	//TODO mark required flags
	tokenCmd.PersistentFlags().StringVarP(&keystoneEndpoint, "keystone-endpoint", "", "", "Keystone endpoint")
	tokenCmd.PersistentFlags().StringVarP(&user, "user", "", "", "Username")
	tokenCmd.PersistentFlags().StringVarP(&password, "password", "", "", "Password")
	tokenCmd.PersistentFlags().StringVarP(&userDomainName, "user-domain-name", "", "", "User Domain Name")
	tokenCmd.PersistentFlags().StringVarP(&project, "project", "", "", "Project")
	tokenCmd.PersistentFlags().StringVarP(&projectDomainName, "project-domain-name", "", "", "Project Domain Name")
	tokenCmd.PersistentFlags().StringVarP(&format, "format", "", "", "text, json, curlrc (Default: text)")

	rootCmd.AddCommand(tokenCmd)

}

func Token(cmd *cobra.Command, args []string) {

	keystoneEndpoint, err := cmd.Flags().GetString("keystone-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	user, err := cmd.Flags().GetString("user")
	if err != nil {
		log.Fatal(err)
	}
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Fatal(err)
	}
	userDomainName, err := cmd.Flags().GetString("user-domain-name")
	if err != nil {
		log.Fatal(err)
	}
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		log.Fatal(err)
	}

	projectDomainName, err := cmd.Flags().GetString("project-domain-name")
	if err != nil {
		log.Fatal(err)
	}

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		log.Fatal(err)
	}

	legacy.Get(userDomainName, project, projectDomainName, user, keystoneEndpoint, password, format)
}
