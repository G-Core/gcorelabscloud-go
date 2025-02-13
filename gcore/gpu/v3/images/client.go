package images

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type ServiceClient struct {
	*gcorecloud.ServiceClient
}

// ImageOpts represents common options for uploading GPU images.
type ImageOpts struct {
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

// ToImageCreateMap builds a request body from ImageOpts.
func (opts ImageOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

func (c *ServiceClient) UploadImage(opts ImageOpts) (*tasks.TaskResults, error) {
	url := ImagesURL(c.ServiceClient)
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

// NewImageOpts creates a new ImageOpts instance
func NewImageOpts() ImageOpts {
	return ImageOpts{}
}

// ListBaremetalImages retrieves list of baremetal GPU images
func (c *ServiceClient) ListBaremetalImages() pagination.Pager {
	return pagination.NewPager(c.ServiceClient, ImagesURL(c.ServiceClient), func(r pagination.PageResult) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListVirtualImages retrieves list of virtual GPU images
func (c *ServiceClient) ListVirtualImages() pagination.Pager {
	return pagination.NewPager(c.ServiceClient, ImagesURL(c.ServiceClient), func(r pagination.PageResult) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}
