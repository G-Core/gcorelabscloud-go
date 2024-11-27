package instances

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v2/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

// ActionOptsBuilder allows extensions to add parameters to the action request.
type ActionOptsBuilder interface {
	ToActionMap() (map[string]interface{}, error)
}

// ActionOpts represents options used to run action.
type ActionOpts struct {
	Action          types.InstanceActionType `json:"action" required:"true" validate:"required,enum"`
	ActivateProfile *bool                    `json:"activate_profile,omitempty"`
}

// Validate checks if the ActionOpts is valid.
func (opts ActionOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToActionMap builds a request body from ActionOpts.
func (opts ActionOpts) ToActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// Action run an action for the instance.
func Action(client *gcorecloud.ServiceClient, instanceID string, opts ActionOptsBuilder) (r tasks.Result) {
	b, err := opts.ToActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(resourceActionURL(client, instanceID), b, &r.Body, nil) // nolint
	return
}

// MetadataItemOpts represents options used to find a metadata.
type MetadataItemOpts struct {
	Key string `q:"key" validate:"required"`
}

type MetadataItemBuilder interface {
	ToMetadataItemQuery() (string, error)
}

// ToMetadataItemQuery builds a query string from MetadataItemOpts.
func (opts MetadataItemOpts) ToMetadataItemQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// MetadataItemGet gets defined metadata key for an instance.
func MetadataItemGet(client *gcorecloud.ServiceClient, id string, opts MetadataItemOpts) (r metadata.MetadataResult) {
	url := metadataItemURL(client, id)
	query, err := opts.ToMetadataItemQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query
	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}

// MetadataItemDelete deletes defined metadata key for an instance.
func MetadataItemDelete(client *gcorecloud.ServiceClient, id string, opts MetadataItemOpts) (r metadata.MetadataActionResult) {
	url := metadataItemURL(client, id)
	query, err := opts.ToMetadataItemQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query
	_, r.Err = client.Delete(url, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}
