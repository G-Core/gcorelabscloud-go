package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const TokenID = "token_id"

var (
	TestUploadBaremetalImageResponse = `
{
    "id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
    "state": "NEW",
    "task_type": "upload_baremetal_image"
}`

	TestUploadVirtualImageResponse = `
{
    "id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
    "state": "NEW",
    "task_type": "upload_virtual_image"
}`

	ExpectedTaskResults = &tasks.Task{
		ID:       "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
		State:    tasks.TaskStateNew,
		TaskType: "upload_baremetal_image",
	}
)
