package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/credentials"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "project_id": 1,
      "name": "docker-io",
      "username": "username",
      "password": "password",
      "registry_url": "registry.example.com"
    }
  ]
}
`

const CreateRequest = `
    {
      "name": "docker-io",
      "username": "username",
      "password": "password",
      "registry_url": "registry.example.com"
    }
`

const UpdateRequest = `
    {
      "username": "username",
      "password": "password",
      "registry_url": "registry.example.com"
    }
`

const GetResponse = `
    {
      "project_id": 1,
      "name": "docker-io",
      "username": "username",
      "password": "password",
      "registry_url": "registry.example.com"
    }
`

var (
	Creds1 = credentials.RegistryCredentials{
		ProjectID:   fake.ProjectID,
		Name:        "docker-io",
		Username:    "username",
		RegistryURL: "registry.example.com",
	}
	Creds1Password = "password"
	CredsSlice     = []credentials.RegistryCredentials{Creds1}
)
