package auth

import (
	"errors"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

type TokenInput struct {
	Region string

	UserID     string
	Username   string
	DomainID   string
	DomainName string
	Password   string

	ProjectID         string
	TenantID          string
	ProjectName       string
	ProjectDomainName string
}

func (t *TokenInput) Validate() error {

	// Region is required
	if t.Region == "" {
		return errors.New("empty region is not supported")
	}

	// If UserID is not provided (username AND (domain id OR domain name)) is required
	if t.UserID != "" {

		if t.Username == "" {
			return errors.New("empty username is not supported when userId is empty")
		} else if t.DomainName == "" && t.DomainID == "" {
			return errors.New("either domain name or domainId required")
		}
	}

	if t.Password == "" {
		return errors.New("empty password is not supported")
	}

	/* Scope is optional
	if t.ProjectID == "" && t.TenantID == "" {
		if t.ProjectName == "" || t.ProjectDomainName == "" {
			return errors.New("project name and project domain name required when project id or tenant id is empty")
		}
	}
	*/

	return nil
}

type TokenProvider interface {
	Get(input TokenInput) (string, error)
}

type OpenStackTokenProvider struct{}

func (o OpenStackTokenProvider) Get(input TokenInput) (string, error) {

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", input.Region),

		UserID:     input.UserID,
		Username:   input.Username,
		DomainID:   input.DomainID,
		DomainName: input.DomainName,
		Password:   input.Password,

		TenantID: input.TenantID, //V2

	}

	// Check for scope
	if input.ProjectID != "" || input.ProjectName == "" || input.ProjectDomainName == "" {
		opts.Scope = &gophercloud.AuthScope{
			ProjectID:   input.ProjectID,
			ProjectName: input.ProjectName,
			DomainName:  input.ProjectDomainName,
		}
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return "", err

	}

	return client.Token(), nil
}
