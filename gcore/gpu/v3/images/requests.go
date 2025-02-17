package images

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

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

// UploadImage uploads a new GPU image
func UploadImage(client *gcorecloud.ServiceClient, opts ImageOpts) (r tasks.Result) {
	url := ImagesURL(client)
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	return
}

// List retrieves list of GPU images
func List(client *gcorecloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, ImagesURL(client), func(r pagination.PageResult) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete deletes a GPU image by ID
func Delete(client *gcorecloud.ServiceClient, imageID string) (r tasks.Result) {
	url := ImageURL(client, imageID)
	_, r.Err = client.Delete(url, &gcorecloud.RequestOpts{
		OkCodes:      []int{200, 201, 202, 204},
		JSONResponse: &r.Body,
	})
	return
}

// Get retrieves a specific GPU image by ID
func Get(client *gcorecloud.ServiceClient, imageID string) (r GetResult) {
	url := ImageURL(client, imageID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}
