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

var User1 = users.User{UserID: 1}
var Token1 = users.ApiToken{Token: "1"}
