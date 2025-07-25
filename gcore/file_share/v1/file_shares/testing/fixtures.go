package testing

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/file_share/v1/file_shares"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

var ListResponse = `
{
  "count": 1,
  "results": [
	{
		"task_id": null,
		"connection_point": "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
		"status": "available",
		"type_name": "default_share_type",
		"creator_task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		"created_at": "2023-08-01T14:32:41.465031",
		"name": "myshare",
		"size": 13,
		"protocol": "NFS",
		"id": "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
		"region": "ED-10",
		"region_id": 2,
		"project_id": 5,
		"metadata_detailed": [
			{
				"key": "qqq",
				"value": "that",
				"read_only": false
			},
			{
				"key": "task_id",
				"value": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
				"read_only": true
			}
		],
		"metadata": {
			"qqq": "that",
			"task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"
		}
	}
  ]
}
`

var GetResponse = `
{
    "task_id": null,
    "connection_point": "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
    "status": "available",
    "type_name": "default_share_type",
    "creator_task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
    "created_at": "2023-08-01T14:32:41.465031",
    "network_name": "usernet",
    "name": "myshare",
    "size": 13,
    "share_network_name": "File_share_ivandshare2_network",
    "protocol": "NFS",
    "id": "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
    "subnet_name": "usersnet",
    "metadata": {
        "task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
        "qqq": "that"
    },
    "region": "ED-10",
    "region_id": 2,
    "project_id": 5,
    "metadata_detailed": [
        {
            "key": "qqq",
            "value": "that",
            "read_only": false
        },
        {
            "key": "task_id",
            "value": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
            "read_only": true
        }
    ]
}
`

var UpdateResponse = `
{
    "task_id": null,
    "connection_point": "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
    "status": "available",
    "type_name": "default_share_type",
    "creator_task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
    "created_at": "2023-08-01T14:32:41.465031",
    "network_name": "usernet",
    "name": "myshareqqq",
    "size": 13,
    "share_network_name": "File_share_ivandshare2_network",
    "protocol": "NFS",
    "id": "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
    "subnet_name": "usersnet",
    "metadata": {
        "task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
        "qqq": "that"
    },
    "region": "ED-10",
    "region_id": 2,
    "project_id": 5,
    "metadata_detailed": [
        {
            "key": "qqq",
            "value": "that",
            "read_only": false
        },
        {
            "key": "task_id",
            "value": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
            "read_only": true
        }
    ]
}
`

const CreateRequest = `
{
    "name": "myshare",
    "network": {
        "network_id": "9b17dd07-1281-4fe0-8c13-d80c5725e297",
        "subnet_id": "221f8318-cf2d-47a7-90f7-97acfa4ef165"
    },
    "protocol": "NFS",
    "size": 13,
    "tags": {
        "qqq": "that"
    }
}
`
const UpdateRequest = `
{
	"name": "myshareqqq"
}
`

const ExtendRequest = `
{
	"size": 15
}
`

const CreateResponse = `
{
  "tasks": [
    "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"
  ]
}
`

const DeleteResponse = `
{
  "tasks": [
    "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"
  ]
}
`

const ExtendResponse = `
{
  "tasks": [
    "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"
  ]
}
`

const CreateAccessRuleRequest = `
{
    "access_mode": "rw",
    "ip_address": "10.100.100.0/24"
}
`

var ListAccessRuleResponse = `
{
    "count": 1,
    "results": [
        {
            "state": "active",
            "access_level": "ro",
            "access_to": "10.17.18.10",
            "id": "6a0a0be1-5875-4a0a-82dd-bab2eef8cb3f"
        }
    ]
}
`

const CreateAccessRuleResponse = `
{
	"state": "active",
	"access_level": "rw",
	"access_to": "10.100.100.0/24",
	"id": "6a0a0be1-5875-4a0a-82dd-bab2eef8cbaa"
}
`

const MetadataResponse = `
{
	"key": "task_id",
	"value": "47e8d97b-0318-4d3e-91c6-53a0e0016f81",
	"read_only": true
}
`
const MetadataCreateRequest = `
{
"test1": "test1",
"test2": "test2"
}
`
const MetadataListResponse = `
{
    "count": 2,
    "results": [
        {
            "key": "qqq",
            "value": "that",
            "read_only": false
        },
        {
            "key": "task_id",
            "value": "47e8d97b-0318-4d3e-91c6-53a0e0016f81",
            "read_only": true
        }
    ]
}
`

var fileShareNetworkConfigWithoutSubnet = `
{
	"network_id": "9b17dd07-1281-4fe0-8c13-d80c5725e297"
}
`

var createdTimeString = "2023-08-01T14:32:41.465031"
var createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339MilliNoZ, createdTimeString)
var createdTime = gcorecloud.JSONRFC3339MilliNoZ{Time: createdTimeParsed}

var creatorTaskID = "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"

var (
	networkName = "usernet"
	subnetName  = "usersnet"
	FileShare1  = file_shares.FileShare{
		Name:             "myshare",
		ID:               "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
		Protocol:         "NFS",
		Status:           file_shares.StatusAvailable,
		Size:             13,
		TypeName:         "cephfs",
		CreatedAt:        &createdTime,
		ShareNetworkName: "File_share_ivandshare2_network",
		NetworkName:      &networkName,
		SubnetName:       &subnetName,
		ConnectionPoint:  "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
		TaskID:           nil,
		CreatorTaskID:    &creatorTaskID,
		ProjectID:        5,
		RegionID:         2,
		Region:           "ED-10",
		Metadata: map[string]interface{}{
			"qqq":     "that",
			"task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		},
	}

	ListFileShare1 = file_shares.FileShare{
		Name:             "myshare",
		ID:               "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
		Protocol:         "NFS",
		Status:           file_shares.StatusAvailable,
		Size:             13,
		TypeName:         "cephfs",
		CreatedAt:        &createdTime,
		ShareNetworkName: "",
		ConnectionPoint:  "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
		TaskID:           nil,
		CreatorTaskID:    &creatorTaskID,
		ProjectID:        5,
		RegionID:         2,
		Region:           "ED-10",
		Metadata: map[string]interface{}{
			"qqq":     "that",
			"task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		},
	}

	AccessRule1 = file_shares.AccessRule{
		ID:          "6a0a0be1-5875-4a0a-82dd-bab2eef8cb3f",
		State:       "active",
		AccessTo:    "10.17.18.10",
		AccessLevel: "ro",
	}

	CreatedAccessRule = file_shares.AccessRule{
		ID:          "6a0a0be1-5875-4a0a-82dd-bab2eef8cbaa",
		State:       "active",
		AccessTo:    "10.100.100.0/24",
		AccessLevel: "rw",
	}

	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"},
	}

	ExpectedFileShareSlice = []file_shares.FileShare{ListFileShare1}

	ExpectedAccessRuleSlice = []file_shares.AccessRule{AccessRule1}

	ResourceMetadata = map[string]interface{}{
		"some_key": "some_val",
	}

	Metadata1 = metadata.Metadata{
		Key:      "qqq",
		Value:    "that",
		ReadOnly: false,
	}
	Metadata2 = metadata.Metadata{
		Key:      "task_id",
		Value:    "47e8d97b-0318-4d3e-91c6-53a0e0016f81",
		ReadOnly: true,
	}
	ExpectedMetadataList = []metadata.Metadata{Metadata1, Metadata2}

	FirstFileShare = file_shares.FileShare{
		TaskID:          nil,
		ConnectionPoint: "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
		Status:          "available",
		TypeName:        "default_share_type",
		CreatorTaskID:   &creatorTaskID,
		CreatedAt:       &createdTime,
		Name:            "myshare",
		Size:            13,
		Protocol:        "NFS",
		ID:              "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
		Region:          "ED-10",
		RegionID:        2,
		ProjectID:       5,
		Metadata:        map[string]interface{}{"task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54", "qqq": "that"},
	}

	SecondFileShare = file_shares.FileShare{
		TaskID:           nil,
		ConnectionPoint:  "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
		Status:           "available",
		TypeName:         "default_share_type",
		CreatorTaskID:    &creatorTaskID,
		CreatedAt:        &createdTime,
		NetworkName:      &networkName,
		Name:             "myshare",
		Size:             13,
		ShareNetworkName: "File_share_ivandshare2_network",
		Protocol:         "NFS",
		ID:               "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
		SubnetName:       &subnetName,
		Metadata:         map[string]interface{}{"task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54", "qqq": "that"},
		Region:           "ED-10",
		RegionID:         2,
		ProjectID:        5,
	}

	UpdatedFileShare = file_shares.FileShare{
		TaskID:           nil,
		ConnectionPoint:  "10.33.20.241:/shares/share-e1dca5e4-257d-47c2-82ac-980fa43e0da9",
		Status:           "available",
		TypeName:         "default_share_type",
		CreatorTaskID:    &creatorTaskID,
		CreatedAt:        &createdTime,
		NetworkName:      &networkName,
		Name:             "myshareqqq",
		Size:             13,
		ShareNetworkName: "File_share_ivandshare2_network",
		Protocol:         "NFS",
		ID:               "8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
		SubnetName:       &subnetName,
		Metadata:         map[string]interface{}{"task_id": "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54", "qqq": "that"},
		Region:           "ED-10",
		RegionID:         2,
		ProjectID:        5,
	}

	FirstAccessRule = file_shares.AccessRule{
		ID:          "6a0a0be1-5875-4a0a-82dd-bab2eef8cb3f",
		State:       "active",
		AccessTo:    "10.17.18.10",
		AccessLevel: "ro",
	}

	SecondAccessRule = file_shares.AccessRule{
		ID:          "6a0a0be1-5875-4a0a-82dd-bab2eef8cbaa",
		State:       "active",
		AccessTo:    "10.100.100.0/24",
		AccessLevel: "rw",
	}

	FirstMetadata = metadata.Metadata{
		Key:      "qqq",
		Value:    "that",
		ReadOnly: false,
	}

	SecondMetadata = metadata.Metadata{
		Key:      "task_id",
		Value:    "47e8d97b-0318-4d3e-91c6-53a0e0016f81",
		ReadOnly: true,
	}
	Task = tasks.Task{
		ID:        "79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		State:     "RUNNING",
		TaskType:  "SFS_CREATE_SHARE",
		CreatedOn: gcorecloud.JSONRFC3339NoZ{Time: createdTime.Time},
		UpdatedOn: &gcorecloud.JSONRFC3339NoZ{Time: createdTime.Time},
		CreatedResources: &map[string]interface{}{
			"file_shares": []interface{}{
				"8fba32f8-dc70-4ac2-be9c-ed6b02927c0e",
			},
		},
	}
)

const CheckLimitsRequest = `
{
    "size": 1001
}
`

const CheckLimitsResponse = `
{
    "sfs_count_limit": 10,
    "sfs_count_requested": 1,
    "sfs_count_usage": 2,
    "sfs_size_limit": 1000,
    "sfs_size_requested": 1001,
    "sfs_size_usage": 900
}
`

var ExpectedCheckLimitsQuota = file_shares.FileShareQuota{
	CountLimit:     10,
	CountRequested: 1,
	CountUsage:     2,
	SizeLimit:      1000,
	SizeRequested:  1001,
	SizeUsage:      900,
}
