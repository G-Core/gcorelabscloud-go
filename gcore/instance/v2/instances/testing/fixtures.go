package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

const MetadataResponse = `
{
  "key": "db.name",
  "value": "pg",
  "read_only": false
}
`

const ActionResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	instanceID = "ad1bb86e-2f83-4e0f-87c0-e1fd777d6352"
	Metadata   = metadata.Metadata{
		Key:      "db.name",
		Value:    "pg",
		ReadOnly: false,
	}
	ActionTasks = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
)
