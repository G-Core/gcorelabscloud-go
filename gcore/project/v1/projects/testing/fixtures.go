package testing

import (
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/types"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/projects"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "id": 1,
      "state": "ACTIVE",
      "created_at": "2020-04-10T11:37:57",
      "description": "",
      "client_id": 1,
      "name": "default"
    }
  ]
}
`

const GetResponse = `
{
  "id": 1,
  "state": "ACTIVE",
  "created_at": "2020-04-10T11:37:57",
  "description": "",
  "client_id": 1,
  "name": "default"
}
`

const CreateRequest = `
{
	"client_id": 1,
	"state": "ACTIVE",
	"name": "default"
}
`

const UpdateRequest = `
{
	"name": "default",
	"description": "description"
}	
`

const CreateResponse = GetResponse
const UpdateResponse = GetResponse
const DeleteResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	createdTimeString    = "2020-04-10T11:37:57"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339NoZ, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339NoZ{Time: createdTimeParsed}

	Project1 = projects.Project{
		ID:          1,
		ClientID:    1,
		Name:        "default",
		Description: "",
		State:       types.ProjectStateActive,
		TaskID:      nil,
		CreatedAt:   createdTime,
	}

	ExpectedProjectSlice = []projects.Project{Project1}
)
