package volumes

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

// InstanceOperationOptsBuilder prepare data to proceed with Attach and Detach requests
type InstanceOperationOptsBuilder interface {
	ToVolumeInstanceOperationMap() (map[string]interface{}, error)
}

// InstanceOperationOpts allows prepare data for Attach and Detach requests
type InstanceOperationOpts struct {
	InstanceID string `json:"instance_id" required:"true" validate:"required,uuid4"`
}

// ToVolumeInstanceOperationMap builds a request body.
func (opts InstanceOperationOpts) ToVolumeInstanceOperationMap() (map[string]interface{}, error) {
	if err := gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts)); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Attach accepts a InstanceOperationOpts struct and attach volume to an instance.
func Attach(c *gcorecloud.ServiceClient, volumeID string, opts InstanceOperationOptsBuilder) (r tasks.Result) {
	b, err := opts.ToVolumeInstanceOperationMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(attachURL(c, volumeID), b, &r.Body, nil)
	return
}

// Detach accepts a InstanceOperationOpts struct and detach volume to an instance.
func Detach(c *gcorecloud.ServiceClient, volumeID string, opts InstanceOperationOptsBuilder) (r tasks.Result) {
	b, err := opts.ToVolumeInstanceOperationMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(detachURL(c, volumeID), b, &r.Body, nil)
	return
}
