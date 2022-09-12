package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/users/v1/users"

const CreateRequest = `
{
  "email": "test@test.test",
  "password": "test"
}
`

const CreateResponse = `
{
  "user_id": 1
}
`

const CreateApiTokenRequest = `
{
  "email": "test@test.test",
  "password": "test",
  "token_name": "test",
  "token_description": "test description"
}
`

const CreateApiTokenResponse = `
{
  "token": "1"
}
`

const UserAssignmentsRequest = `
{
    "client_id": 8,
    "project_id": null,
    "role": "ClientAdministrator",
    "user_id": 777
}
`

const UserAssignmentsResponse = `
{
    "client_id": 8,
    "id": 12,
    "project_id": null,
    "role": "ClientAdministrator",
    "user_id": 777
}
`

var (
	User1  = users.User{UserID: 1}
	Token1 = users.ApiToken{Token: "1"}

	clientID = 8
	UA1      = users.UserAssignment{
		ClientID: &clientID,
		ID:       12,
		Role:     "ClientAdministrator",
		UserID:   777,
	}
)
