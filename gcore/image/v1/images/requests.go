package images

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images/types"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToImageListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Private    bool              `q:"private" validate:"omitempty"`
	Visibility types.Visibility  `q:"visibility" validate:"omitempty"`
	MetadataK  string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV map[string]string `q:"metadata_kv" validate:"omitempty"`
}

// ToImageListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToImageListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToImageCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create an image.
type CreateOpts struct {
	Name           string                      `json:"name" required:"true" validate:"required"`
	HwMachineType  types.HwMachineType         `json:"hw_machine_type,omitempty" validate:"enum"`
	SshKey         types.SshKeyType            `json:"ssh_key,omitempty" validate:"required,enum"`
	OSType         types.OSType                `json:"os_type" validate:"required,enum"`
	IsBaremetal    *bool                       `json:"is_baremetal,omitempty"`
	HwFirmwareType types.HwFirmwareType        `json:"hw_firmware_type,omitempty" validate:"enum"`
	Source         types.ImageSourceType       `json:"source" validate:"required,enum"`
	VolumeID       string                      `json:"volume_id" required:"true" validate:"required"`
	Metadata       map[string]string           `json:"metadata,omitempty"`
	Architecture   types.ImageArchitectureType `json:"architecture,omitempty" validate:"enum"`
}

/*
       "cow_format": False,
       "hw_firmware_type": "bios",
       "hw_machine_type": "q35",
       "is_baremetal": False,
       "name": "image_name",
       "os_type": "linux",
       "ssh_key": "allow",
       "url": "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img",
       "metadata": {"key": "value"},
   }
*/
// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToImageCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToImageUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to create an image.
type UpdateOpts struct {
	HwMachineType  types.HwMachineType  `json:"hw_machine_type" validate:"required,enum"`
	SshKey         types.SshKeyType     `json:"ssh_key" validate:"required,enum"`
	Name           string               `json:"name" required:"true"`
	OSType         types.OSType         `json:"os_type" validate:"required,enum"`
	IsBaremetal    *bool                `json:"is_baremetal,omitempty"`
	HwFirmwareType types.HwFirmwareType `json:"hw_firmware_type" validate:"required,enum"`
}

// Validate
func (opts UpdateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToImageUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToImageUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// UploadOptsBuilder allows extensions to add additional parameters to the Upload request.
type UploadOptsBuilder interface {
	ToImageUploadMap() (map[string]interface{}, error)
}

// UploadOpts represents options used to upload an image.
type UploadOpts struct {
	OsVersion      string                      `json:"os_version,omitempty"`
	HwMachineType  types.HwMachineType         `json:"hw_machine_type,omitempty" validate:"enum"`
	SshKey         types.SshKeyType            `json:"ssh_key,omitempty" validate:"enum"`
	Name           string                      `json:"name" required:"true" validate:"required"`
	OsDistro       string                      `json:"os_distro,omitempty"`
	OSType         types.OSType                `json:"os_type" validate:"enum"`
	URL            string                      `json:"url" required:"true" validate:"required,url"`
	IsBaremetal    *bool                       `json:"is_baremetal,omitempty"`
	HwFirmwareType types.HwFirmwareType        `json:"hw_firmware_type,omitempty" validate:"enum"`
	CowFormat      bool                        `json:"cow_format"`
	Metadata       map[string]string           `json:"metadata,omitempty"`
	Architecture   types.ImageArchitectureType `json:"architecture,omitempty" validate:"enum"`
}

// Validate
func (opts UploadOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToImageUploadMap builds a request body from UploadOpts.
func (opts UploadOpts) ToImageUploadMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

func List(client *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToImageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific image based on its unique ID.
func Get(client *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(client, id)
	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}

func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Image, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractImages(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// Create an image.
func Create(client *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, nil) // nolint
	return
}

// Delete an image.
func Delete(client *gcorecloud.ServiceClient, imageID string) (r tasks.Result) {
	url := deleteURL(client, imageID)
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
	return
}

// Update accepts a UpdateOpts struct and updates an existing image using the
// values provided.
func Update(client *gcorecloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	url := updateURL(client, id)
	b, err := opts.ToImageUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(url, b, &r.Body, nil) // nolint
	return
}

// Upload accepts a UploadOpts struct and upload an image using the
// values provided.
func Upload(client *gcorecloud.ServiceClient, opts UploadOptsBuilder) (r tasks.Result) {
	b, err := opts.ToImageUploadMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(uploadURL(client), b, &r.Body, nil) // nolint
	return
}
