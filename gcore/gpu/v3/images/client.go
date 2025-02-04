package images

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

type ServiceClient struct {
	*gcorecloud.ServiceClient
}

// UploadBaremetalImageOpts represents options for uploading a baremetal GPU image.
type UploadBaremetalImageOpts struct {
	// Image name
	Name string `json:"name" required:"true"`

	// Image URL
	URL string `json:"url" required:"true"`

	// Image architecture type: aarch64, x86_64
	Architecture *string `json:"architecture,omitempty"`

	// When True, image cannot be deleted unless all volumes, created from it, are deleted.
	CowFormat *bool `json:"cow_format,omitempty"`

	// Specifies the type of firmware with which to boot the guest.
	HwFirmwareType *ImageHwFirmwareType `json:"hw_firmware_type,omitempty"`

	// Create one or more metadata items for a cluster
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// OS Distribution, i.e. Debian, CentOS, Ubuntu, CoreOS etc.
	OsDistro *string `json:"os_distro,omitempty"`

	// The operating system installed on the image.
	OsType *ImageOsType `json:"os_type,omitempty"`

	// OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian
	OsVersion *string `json:"os_version,omitempty"`

	// Permission to use a ssh key in instances
	SshKey *SshKeyType `json:"ssh_key,omitempty"`
}

// UploadVirtualImageOpts represents options for uploading a virtual GPU image.
type UploadVirtualImageOpts struct {
	// Image name
	Name string `json:"name" required:"true"`

	// Image URL
	URL string `json:"url" required:"true"`

	// Image architecture type: aarch64, x86_64
	Architecture *string `json:"architecture,omitempty"`

	// When True, image cannot be deleted unless all volumes, created from it, are deleted.
	CowFormat *bool `json:"cow_format,omitempty"`

	// Specifies the type of firmware with which to boot the guest.
	HwFirmwareType *ImageHwFirmwareType `json:"hw_firmware_type,omitempty"`

	// Create one or more metadata items for a cluster
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// OS Distribution, i.e. Debian, CentOS, Ubuntu, CoreOS etc.
	OsDistro *string `json:"os_distro,omitempty"`

	// The operating system installed on the image.
	OsType *ImageOsType `json:"os_type,omitempty"`

	// OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian
	OsVersion *string `json:"os_version,omitempty"`

	// Permission to use a ssh key in instances
	SshKey *SshKeyType `json:"ssh_key,omitempty"`
}

// ToImageCreateMap builds a request body from UploadBaremetalImageOpts.
func (opts UploadBaremetalImageOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToImageCreateMap builds a request body from UploadVirtualImageOpts.
func (opts UploadVirtualImageOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

func (c *ServiceClient) UploadBaremetalImage(opts UploadBaremetalImageOpts) (*tasks.TaskResults, error) {
	url := UploadBaremetalURL(c.ServiceClient)
	b, err := opts.ToImageCreateMap()
	if err != nil {
		return nil, err
	}
	var result tasks.TaskResults
	_, err = c.Post(url, b, &result, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ServiceClient) UploadVirtualImage(opts UploadVirtualImageOpts) (*tasks.TaskResults, error) {
	url := UploadVirtualURL(c.ServiceClient)
	b, err := opts.ToImageCreateMap()
	if err != nil {
		return nil, err
	}
	var result tasks.TaskResults
	_, err = c.Post(url, b, &result, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
