package testing

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
)

const AttachDetachRequest = `
{
  "instance_id": "8dc30d49-bb34-4920-9bbd-03a2587ec0ad"
}	
`

const AttachResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const DetachResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var rawAttachTask = []byte(`
{
	"job_id": null,
	"request_id": "76f8464efaa7f41bb201e6aaca309262",
	"state": "FINISHED",
	"error": null,
	"id": "c44ba6f3-e371-498f-91bb-618eb38d1b3a",
	"acknowledged_by": null,
	"lifecycle_policy_id": null,
	"detailed_state": null,
	"schedule_id": null,
	"region_id": 76,
	"data": {
		"volume_id": "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		"instance_id": "dab284e9-1fb5-47e8-870e-3b67a350ac17",
		"attachment_tag": null
	},
	"task_type": "attach_volume",
	"client_id": 65393,
	"updated_on": "2024-05-07T08:24:02",
	"acknowledged_at": null,
	"project_id": 302434,
	"created_resources": null,
	"user_id": 18339,
	"user_client_id": 65393,
	"created_on": "2024-05-07T08:23:54",
	"finished_on": "2024-05-07T08:24:02"
}
`)

const createdTimeString = "2019-05-29T05:32:41+0000"
const updatedTimeString = "2019-05-29T05:39:20+0000"
const attachedTimeString = "2019-07-26T14:22:03+0000"

var createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
var createdTime = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
var updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
var updatedTime = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
var attachedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, attachedTimeString)
var attachedTime = gcorecloud.JSONRFC3339Z{Time: attachedTimeParsed}

var (
	Volume1 = volumes.Volume{
		AvailabilityZone: "nova",
		CreatedAt:        createdTime,
		UpdatedAt:        updatedTime,
		VolumeType:       "standard",
		ID:               "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		Name:             "123",
		RegionName:       "Luxembourg 1",
		Status:           "available",
		Size:             2,
		Bootable:         false,
		ProjectID:        1,
		RegionID:         1,
		Attachments: []volumes.Attachment{{
			ServerID:     "8dc30d49-bb34-4920-9bbd-03a2587ec0ad",
			AttachmentID: "f2ed59d9-8068-400c-be4b-c4501ef6f33c",
			InstanceName: "123",
			AttachedAt:   attachedTime,
			VolumeID:     "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
			Device:       "/dev/vda",
		},
		},
		Metadata:      []metadata.Metadata{ResourceMetadataReadOnly},
		CreatorTaskID: "d74c2bb9-cea7-4b23-a009-2f13518ae66d",
		VolumeImageMetadata: volumes.VolumeImageMetadata{
			ContainerFormat:               "bare",
			MinRAM:                        "0",
			OwnerSpecifiedOpenstackSHA256: "87ddf8eea6504b5eb849e418a568c4985d3cea59b5a5d069e1dc644de676b4ec",
			DiskFormat:                    "raw",
			ImageName:                     "cirros-gcloud",
			ImageID:                       "723037e2-ec6d-47eb-92de-6276c8907839",
			OwnerSpecifiedOpenstackObject: "images/cirros-gcloud",
			OwnerSpecifiedOpenstackMD5:    "ba3cd24377dde5dfdd58728894004abb",
			MinDisk:                       "1",
			Checksum:                      "ba3cd24377dde5dfdd58728894004abb",
			Size:                          "46137344",
		},
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ResourceMetadataReadOnly = metadata.Metadata{
		Key:      "some_key",
		Value:    "some_val",
		ReadOnly: false,
	}
)
