package legacy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type AuthToken struct {
	Auth `json:"auth"`
}
type Auth struct {
	Identity `json:"identity"`
	Scope    `json:"scope"`
}
type Identity struct {
	Methods  []string `json:"methods"`
	Password `json:"password"`
}
type Password struct {
	User `json:"user"`
}
type User struct {
	Name     string `json:"name"`
	Domain   `json:"domain"`
	Password string `json:"password"`
}
type Project struct {
	Name   string `json:"name"`
	Domain `json:"domain"`
}
type Domain struct {
	Name string `json:"name"`
}
type Scope struct {
	Project `json:"project"`
}

type Token_wraper struct {
	Token Token `json:"token"`
}
type Token struct {
	Is_domain        bool      `json:"is_domain"`
	Methods          []string  `json:"methods"`
	Roles            []Role    `json:"roles"`
	Is_admin_project bool      `json:"is_admin_project"`
	Project          Project_t `json:"project"`
	Expires_at       string    `json:"expires_at"`
	User             Project_t `json:"user"`
	Audit_ids        []string  `json:"audit_ids"`
	Issued_at        string    `json:"issued_at"`
}
type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type Project_t struct {
	Domain Domain_t `json:"domain"`
	Id     string   `json:"id"`
	Name   string   `json:"name"`
}
type Domain_t struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func Get(user_domain, project, project_domain, username, keystone_endpoint, password, format string) {
	if password == "" { // If password is not provided from CLI or envvar, try to find and set

		passwordOutput, err := exec.Command("security", "find-generic-password", "-a", username, "-s", "openstack", "-w").Output()
		if err != nil {
			log.Printf("find-generic-password error: %v", err) // Error could occur if it is not found, error content should be checked
		}
		if last := len(passwordOutput) - 1; last >= 0 && passwordOutput[last] == '\n' {
			password = string(passwordOutput[:last])
		}
		if len(passwordOutput) < 2 {
			fmt.Printf("Enter password for user %s: ", username)
			fmt.Scanln(&password)
		}
	}

	payload, err := payload(username, user_domain, password, project, project_domain)
	if err != nil {
		log.Fatal(err)
	}

	resp, body, obj, err := retrieve(keystone_endpoint, payload)
	if err != nil {
		log.Fatal(err)
	}
	// output handling
	print(format, body, resp, obj)

}
func print(format string, body []byte, resp *http.Response, obj Token_wraper) {
	if format == "json" {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, body, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(prettyJSON.Bytes()))
	} else if format == "curlrc" {
		fmt.Println("header \"X-Auth-Token: ", resp.Header.Get("X-Subject-Token"), "\"")
		fmt.Println("header \"Content-Type: application/json\"")
	} else {
		fmt.Println(resp.Header.Get("X-Subject-Token"))
		fmt.Println("User:\t\t", obj.Token.User.Id)
		fmt.Println("Project:\t", obj.Token.Project.Id)
		fmt.Println("Project domain:\t", obj.Token.Project.Domain.Id)
		fmt.Print("Roles:\t\t")
		for _, r := range obj.Token.Roles {
			fmt.Print(r.Name, ", ")
		}
		fmt.Println()
	}
}
func retrieve(keystone_endpoint string, payloadJson []byte) (*http.Response, []byte, Token_wraper, error) {

	url := keystone_endpoint + "/auth/tokens?nocatalog"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJson))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var obj Token_wraper
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		log.Fatal(err)
	}
	return resp, body, obj, nil
}
func payload(username string, user_domain string, password string, project string, project_domain string) ([]byte, error) {
	payload := AuthToken{
		Auth: Auth{
			Identity: Identity{
				Methods: []string{"password"},
				Password: Password{
					User: User{
						Name: username,
						Domain: Domain{
							Name: user_domain,
						},
						Password: string(password),
					},
				},
			},
			Scope: Scope{
				Project: Project{
					Name: project,
					Domain: Domain{
						Name: project_domain,
					},
				},
			},
		},
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return payloadJson, nil
}
