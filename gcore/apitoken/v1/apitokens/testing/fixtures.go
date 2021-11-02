package testing

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/apitokens"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/types"
	"time"
)

const ListResponse = `
[
  {
    "name": "My token",
    "description": "It's my token",
    "exp_date": null,
    "client_user": {
      "role": {
        "id": 1,
        "name": "Administrators"
      },
      "deleted": false,
      "user_id": 123,
      "user_name": "John Doe",
      "user_email": "some@email.com",
      "client_id": 456
    },
    "id": 42,
    "deleted": false,
    "expired": false,
    "created": "2021-01-01T12:00:00.000Z",
    "last_usage": null
  }
]
`

const GetResponse = `
{
  "name": "My token",
  "description": "It's my token",
  "exp_date": null,
  "client_user": {
    "role": {
      "id": 1,
      "name": "Administrators"
    },
    "deleted": false,
    "user_id": 123,
    "user_name": "John Doe",
    "user_email": "some@email.com",
    "client_id": 456
  },
  "id": 42,
  "deleted": false,
  "expired": false,
  "created": "2021-01-01T12:00:00.000Z",
  "last_usage": null
}
`

const CreateRequest = `
{
  "name": "My token",
  "description": "It's my token",
  "exp_date": null,
  "client_user": {
    "role": {
      "id": 1,
      "name": "Administrators"
    }
  }
}
`

const CreateResponse = `
{
  "token": "string"
}
`

var (
	clientID             = 3
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339ZZ, "2021-01-01T12:00:00.000Z")
	createdTime          = gcorecloud.JSONRFC3339ZZ{Time: createdTimeParsed}

	apiToken1 = apitokens.APIToken{
		ID:          42,
		Name:        "My token",
		Description: "It's my token",
		ExpDate:     nil,
		ClientUser: &apitokens.ClientUser{
			Role: apitokens.ClientRole{
				ID:   types.RoleIDAdministrators,
				Name: types.RoleNameAdministrators,
			},
			Deleted:   false,
			UserID:    123,
			UserName:  "John Doe",
			UserEmail: "some@email.com",
			ClientID:  456,
		},
		Deleted:   false,
		Expired:   false,
		Created:   createdTime,
		LastUsage: nil,
	}
	ExpectedAPITokenSlice = []apitokens.APIToken{apiToken1}
	Token1                = apitokens.Token{Token: "string"}
)
