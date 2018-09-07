package command

import (
	"fmt"
	"github.com/sapcc/token-tool/pkg/auth"
	"github.com/spf13/cobra"
	"log"
)

var input auth.TokenInput
var provider auth.TokenProvider

func init() {

	var tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "Retrieves token from Keystone",
		Long:  ``,
		Run:   Token,
	}

	provider = auth.OpenStackTokenProvider{}

	tokenCmd.PersistentFlags().StringVarP(&input.Region, "region", "", "", "Region")
	tokenCmd.PersistentFlags().StringVarP(&input.UserID, "user-id", "", "", "User ID")
	tokenCmd.PersistentFlags().StringVarP(&input.Username, "username", "", "", "Username")
	tokenCmd.PersistentFlags().StringVarP(&input.DomainID, "domain-id", "", "", "Domain ID")
	tokenCmd.PersistentFlags().StringVarP(&input.DomainName, "domain-name", "", "", "Domain Name")
	tokenCmd.PersistentFlags().StringVarP(&input.Password, "password", "", "", "Password")

	tokenCmd.PersistentFlags().StringVarP(&input.ProjectID, "project-id", "", "", "Project ID")
	tokenCmd.PersistentFlags().StringVarP(&input.TenantID, "tenant-id", "", "", "Tenant ID")
	tokenCmd.PersistentFlags().StringVarP(&input.ProjectName, "project-name", "", "", "Project Name")
	tokenCmd.PersistentFlags().StringVarP(&input.ProjectDomainName, "project-domain-name", "", "", "Project Domain Name")

	rootCmd.AddCommand(tokenCmd)

}

func Token(cmd *cobra.Command, args []string) {

	err := input.Validate()
	if err != nil {
		log.Fatal(err)
	}

	token, err := provider.Get(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)
}
