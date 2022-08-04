package testing

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

var ListResponse = fmt.Sprintf(`
{
  "count": 1,
  "results": [
   {
      "id": "4a44e5a2-e7ba-41b8-bf78-ddfa2e22974b",
      "project_id": 1,
      "min_ram": 0,
      "region_id": 1,
      "visibility": "public",
      "os_distro": "fedora-coreos",
      "tags": [],
      "updated_at": "2020-03-09T10:16:54+0000",
      "size": 1685454848,
      "task_id": null,
      "region": "RegionOne",
      "created_at": "2020-03-09T10:16:45+0000",
      "disk_format": "qcow2",
      "min_disk": 0,
      "name": "fedora-coreos",
      "status": "active",
      "metadata_detailed": [%s],
  	  "metadata": %s
    }
  ]
}
`, MetadataResponse, MetadataResponseKV)

const MetadataResponse = `
{
  "key": "some_key",
  "value": "some_val",
  "read_only": false
}
`
const MetadataResponseKV = `
{
  "some_key": "some_val"
}
`

var GetResponse = fmt.Sprintf(`
{
  "id": "4a44e5a2-e7ba-41b8-bf78-ddfa2e22974b",
  "project_id": 1,
  "min_ram": 0,
  "region_id": 1,
  "visibility": "public",
  "os_distro": "fedora-coreos",
  "tags": [],
  "updated_at": "2020-03-09T10:16:54+0000",
  "size": 1685454848,
  "task_id": null,
  "region": "RegionOne",
  "created_at": "2020-03-09T10:16:45+0000",
  "disk_format": "qcow2",
  "min_disk": 0,
  "name": "fedora-coreos",
  "status": "active",
  "metadata_detailed": [%s],
  "metadata": %s
}
`, MetadataResponse, MetadataResponseKV)

const CreateRequest = `
{
  "hw_firmware_type": "bios",
  "hw_machine_type": "q35",
  "is_baremetal": false,
  "name": "test_image",
  "os_type": "linux",
  "source": "volume",
  "ssh_key": "allow",
  "volume_id": "d478ae29-dedc-4869-82f0-96104425f565"
}
`

const UploadRequest = `
{
  "cow_format": false,
  "hw_firmware_type": "bios",
  "hw_machine_type": "q35",
  "is_baremetal": false,
  "name": "image_name",
  "os_type": "linux",
  "ssh_key": "allow",
  "url": "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"
}
`

const UpdateRequest = `
{
  "hw_machine_type": "i440",
  "ssh_key": "allow",
  "name": "string",
  "os_type": "linux",
  "is_baremetal": true,
  "hw_firmware_type": "bios"
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
	ctm, _      = time.Parse(gcorecloud.RFC3339Z, "2020-03-09T10:16:45+0000")
	createdTime = gcorecloud.JSONRFC3339Z{Time: ctm}
	utm, _      = time.Parse(gcorecloud.RFC3339Z, "2020-03-09T10:16:54+0000")
	updatedTime = gcorecloud.JSONRFC3339Z{Time: utm}

	Image1 = images.Image{
		ID:            "4a44e5a2-e7ba-41b8-bf78-ddfa2e22974b",
		Name:          "fedora-coreos",
		Description:   "",
		Status:        "active",
		Tags:          []string{},
		Visibility:    "public",
		MinDisk:       0,
		MinRAM:        0,
		OsDistro:      "fedora-coreos",
		OsVersion:     "",
		DisplayOrder:  0,
		CreatedAt:     createdTime,
		UpdatedAt:     &updatedTime,
		Size:          1685454848,
		CreatorTaskID: nil,
		TaskID:        nil,
		DiskFormat:    "qcow2",
		Region:        "RegionOne",
		Metadata:      []metadata.Metadata{ResourceMetadataReadOnly},
	}
	ExpectedImagesSlice = []images.Image{Image1}
	Tasks1              = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ResourceMetadata = map[string]interface{}{
		"some_key": "some_val",
	}

	ResourceMetadataReadOnly = metadata.Metadata{
		Key:      "some_key",
		Value:    "some_val",
		ReadOnly: false,
	}

	Metadata1 = metadata.Metadata{
		Key:      "cost-center",
		Value:    "Atlanta",
		ReadOnly: false,
	}
	Metadata2 = metadata.Metadata{
		Key:      "data-center",
		Value:    "A",
		ReadOnly: false,
	}
	ExpectedMetadataList = []metadata.Metadata{Metadata1, Metadata2}
)
