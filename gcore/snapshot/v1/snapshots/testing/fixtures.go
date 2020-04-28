package testing

import (
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/snapshot/v1/snapshots"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "id": "726ecfcc-7fd0-4e30-a86e-7892524aa483",
      "name": "123",
      "status": "available",
      "description": "mysnapshot",
      "created_at": "2019-05-29T05:32:41+0000",
      "updated_at": "2019-05-29T05:39:20+0000",
      "size": 2,
      "creator_task_id": "2358e3b1-5c42-4705-8950-6ddcfc19c3bd",
      "volume_id": "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
      "project_id": 1,
      "region_id": 1,
	  "region": "RegionOne"	
    }
  ]
}
`

const GetResponse = `
{
  "id": "726ecfcc-7fd0-4e30-a86e-7892524aa483",
  "name": "123",
  "status": "available",
  "description": "mysnapshot",
  "created_at": "2019-05-29T05:32:41+0000",
  "updated_at": "2019-05-29T05:39:20+0000",
  "size": 2,
  "creator_task_id": "2358e3b1-5c42-4705-8950-6ddcfc19c3bd",
  "volume_id": "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
  "project_id": 1,
  "region_id": 1,
  "region": "RegionOne"	
}
`

const CreateRequest = `
{
  "volume_id": "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
  "name": "123",
  "description": "after boot"
}
`

const CreateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const DeleteResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	createdTimeString    = "2019-05-29T05:32:41+0000"
	updatedTimeString    = "2019-05-29T05:39:20+0000"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
	updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
	updatedTime          = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
	creatorTaskID        = "2358e3b1-5c42-4705-8950-6ddcfc19c3bd"
	volumeID             = "67baa7d1-08ea-4fc5-bef2-6b2465b7d227"

	Snapshot1 = snapshots.Snapshot{
		ID:            "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		Name:          "123",
		Description:   "mysnapshot",
		Status:        "available",
		Size:          2,
		VolumeID:      volumeID,
		CreatedAt:     createdTime,
		UpdatedAt:     &updatedTime,
		Metadata:      nil,
		CreatorTaskID: &creatorTaskID,
		TaskID:        nil,
		ProjectID:     fake.ProjectID,
		RegionID:      fake.RegionID,
		Region:        "RegionOne",
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedSnapshotSlice = []snapshots.Snapshot{Snapshot1}
)
