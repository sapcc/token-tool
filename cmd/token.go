package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh/terminal"
)

var version string = "HEAD"

func main() {
	var authInfo clientconfig.AuthInfo
	// handling args/flags
	app := cli.NewApp()
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "username",
			Usage:       "Username. Default: $USER",
			EnvVar:      "OS_USERNAME",
			Destination: &authInfo.Username,
		},
		cli.StringFlag{
			Name:        "user-domain-name",
			Usage:       "User domain name",
			EnvVar:      "OS_USER_DOMAIN_NAME",
			Destination: &authInfo.UserDomainName,
		},
		cli.StringFlag{
			Name:        "user-domain-id",
			Usage:       "User domain ID",
			EnvVar:      "OS_USER_DOMAIN_ID",
			Destination: &authInfo.UserDomainID,
		},
		cli.StringFlag{
			Name:        "user-id",
			Usage:       "User ID",
			EnvVar:      "OS_USER_ID",
			Destination: &authInfo.UserID,
		},
		cli.StringFlag{
			Name:        "project-name",
			Usage:       "Project name",
			EnvVar:      "OS_PROJECT_NAME",
			Destination: &authInfo.ProjectName,
		},
		cli.StringFlag{
			Name:        "project-domain-name",
			Usage:       "Project domain name",
			EnvVar:      "OS_PROJECT_DOMAIN_NAME",
			Destination: &authInfo.ProjectDomainName,
		},
		cli.StringFlag{
			Name:        "project-domain-id",
			Usage:       "Project domain ID",
			EnvVar:      "OS_PROJECT_DOMAIN_ID",
			Destination: &authInfo.ProjectDomainName,
		},
		cli.StringFlag{
			Name:        "project-id",
			Usage:       "Project ID",
			EnvVar:      "OS_PROJECT_ID, OS_TENANT_ID",
			Destination: &authInfo.ProjectID,
		},
		cli.StringFlag{
			Name:        "auth-url",
			Usage:       "keystone/identity endpoint URL",
			EnvVar:      "OS_AUTH_URL",
			Destination: &authInfo.AuthURL,
		},
		cli.StringFlag{
			Name:        "application-credential-id",
			Usage:       "Application Credential ID",
			EnvVar:      "OS_APPLICATION_CREDENTIAL_ID",
			Destination: &authInfo.ApplicationCredentialID,
		},
		cli.StringFlag{
			Name:        "application-credential-name",
			Usage:       "Application Credential Name",
			EnvVar:      "OS_APPLICATION_CREDENTIAL_NAME",
			Destination: &authInfo.ApplicationCredentialName,
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "text",
			Usage: "Format: text, json, curlrc",
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	var authOpts *gophercloud.AuthOptions
	app.Before = func(*cli.Context) (err error) {
		if authOpts, err = clientconfig.AuthOptions(&clientconfig.ClientOpts{AuthInfo: &authInfo}); err != nil {
			return
		}
		//default to system user if no user user variable set
		if authOpts.Username == "" && authOpts.UserID == "" && authOpts.ApplicationCredentialID == "" && authOpts.ApplicationCredentialName == "" {
			authOpts.Username = os.Getenv("USER")
		}
		//if no domain information is given for username we default it top the scope domain name/id
		if authOpts.Username != "" && authOpts.DomainName == "" && authOpts.DomainID == "" && authOpts.Scope != nil {
			if authOpts.Scope.DomainID != "" {
				authOpts.DomainID = authOpts.Scope.DomainID
			} else {
				authOpts.DomainName = authOpts.Scope.DomainName
			}
		}
		//try to get password from keyring if not set via env
		if authOpts.Username != "" && authOpts.Password == "" && authOpts.ApplicationCredentialSecret == "" {
			if pw, err := keyring.Get("openstack", authOpts.Username); err == nil {
				log.Println("Using password from keyring")
				authOpts.Password = pw
			} else {
				if terminal.IsTerminal(int(os.Stdin.Fd())) {
					if password, err := gopass.GetPasswdPrompt("Password: ", true, os.Stdin, os.Stderr); err == nil {
						authOpts.Password = string(password)
					}
				} else {
					if in, err := ioutil.ReadAll(os.Stdin); err == nil && len(in) > 0 {
						log.Println("Password read from stdin")
						authOpts.Password = strings.TrimRight(string(in), "\r\n")
					}
				}
			}
		}

		return

	}

	app.Action = func(c *cli.Context) error {
		return tokenCommand(c.String("format"), authOpts)
	}
	app.Commands = []cli.Command{
		{
			Name:  "curl",
			Usage: "use curl with openstack credentials",
			Action: func(c *cli.Context) error {
				return curlCommand(c.Args(), authOpts)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func tokenCommand(format string, authOptions *gophercloud.AuthOptions) error {
	providerClient, err := openstack.AuthenticatedClient(*authOptions)
	if err != nil {
		return fmt.Errorf("Failed to authenticate: %s", err)
	}
	tokenResponse, ok := providerClient.GetAuthResult().(tokens.CreateResult)
	if !ok {
		return errors.New("Auth response is not a v3 response")
	}

	switch format {
	case "curlrc":
		fmt.Printf("header \"X-Auth-Token: %s\"\n", providerClient.Token())
		fmt.Printf("header \"Content-Type: application/json\"\n")
	case "json":
		//add the token from the heder to the nested json as token_id
		b := tokenResponse.Body.(map[string]interface{})
		b["token_id"] = providerClient.Token()
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		return e.Encode(b)
	default:
		fmt.Println(providerClient.Token())
	}

	return nil
}

func curlCommand(curlArgs []string, authOptions *gophercloud.AuthOptions) error {
	curlPath, err := exec.LookPath("curl")
	if err != nil {
		return fmt.Errorf("curl command not found in path: %s", err)
	}

	providerClient, err := openstack.AuthenticatedClient(*authOptions)
	if err != nil {
		return fmt.Errorf("Failed to authenticate: %s", err)
	}
	tokenResponse, ok := providerClient.GetAuthResult().(tokens.CreateResult)
	if !ok {
		return errors.New("Auth response is not a v3 response")
	}

	catalog, err := tokenResponse.ExtractServiceCatalog()
	if err != nil {
		return fmt.Errorf("Failed to get catalog from auth response: %s", err)
	}

	vars := map[string]string{}
	for _, entry := range catalog.Entries {
		for _, ep := range entry.Endpoints {
			vars[strings.ToUpper(entry.Type)+"_"+strings.ToUpper(ep.Interface)] = ep.URL
			if ep.Interface == "public" {
				vars[strings.ToUpper(entry.Type)] = ep.URL
			}
		}
	}
	for i, arg := range curlArgs {
		curlArgs[i] = os.Expand(arg, func(s string) string { return vars[s] })
	}

	log.Println("curl", strings.Join(curlArgs, " "))
	curlArgs = append([]string{
		curlPath,
		"--header", "X-Auth-Token: " + providerClient.Token(),
		"--header", "Content-Type: application/json",
	},
		curlArgs...)

	return syscall.Exec(curlPath, curlArgs, os.Environ())

}
