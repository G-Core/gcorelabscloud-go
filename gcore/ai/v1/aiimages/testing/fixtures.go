package testing

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/aiimages"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

const ListResponse = `
{
    "count": 1,
    "results": [
        {
            "min_ram": 0,
            "task_id": null,
            "disk_format": "qcow2",
            "os_type": "linux",
            "visibility": "public",
            "min_disk": 0,
            "updated_at": "2022-11-17T11:38:03+0000",
            "project_id": 516070,
            "region_id": 7,
            "region": "ED-10 Preprod",
            "architecture": "x86_64",
            "name": "ubuntu-18.04-x64-poplar-ironic-1.18.0-3.0.0",
            "os_distro": "poplar-ubuntu",
            "created_at": "2022-11-17T11:22:55+0000",
            "vipu_version": "1.18.0",
            "size": 5389418496,
            "sdk_version": "3.0.0",
            "id": "f6aa6e75-ab88-4c19-889d-79133366cb83",
            "metadata_detailed": [],
            "status": "active",
            "os_version": "18.04",
            "metadata": {}
        }
	]
}
`

var (
	ImageCreatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2022-11-17T11:22:55+0000")
	ImageUpdatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2022-11-17T11:38:03+0000")
	AIImage1          = aiimages.AIImage{
		ID:           "f6aa6e75-ab88-4c19-889d-79133366cb83",
		Name:         "ubuntu-18.04-x64-poplar-ironic-1.18.0-3.0.0",
		Status:       "active",
		Visibility:   "public",
		MinDisk:      0,
		MinRAM:       0,
		OsDistro:     "poplar-ubuntu",
		OsType:       "linux",
		OsVersion:    "18.04",
		CreatedAt:    gcorecloud.JSONRFC3339Z{Time: ImageCreatedAt},
		UpdatedAt:    &gcorecloud.JSONRFC3339Z{Time: ImageUpdatedAt},
		Size:         5389418496,
		Region:       "ED-10 Preprod",
		RegionID:     7,
		ProjectID:    516070,
		DiskFormat:   "qcow2",
		Architecture: "x86_64",
		Metadata:     []metadata.Metadata{},
	}
	ExpectedAIImageSlice = []aiimages.AIImage{AIImage1}
)
