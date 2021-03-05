package snapshots

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToSnapshotListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	VolumeID   string `q:"volume_id"`
	InstanceID string `q:"instance_id"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSnapshotListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToSnapshotListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SnapshotPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific snapshot based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToSnapshotCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a snapshot.
type CreateOpts struct {
	VolumeID    string `json:"volume_id" required:"true" validate:"required"`
	Name        string `json:"name" required:"true" validate:"required"`
	Description string `json:"description,omitempty"`
}

// ToSnapshotCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSnapshotCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new snapshot using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToSnapshotCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToSnapshotUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a snapshot.
type UpdateOpts struct {
	Name string `json:"name,omitempty"`
}

// ToSnapshotUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToSnapshotUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Delete accepts a unique ID and deletes the snapshot associated with it.
func Delete(c *gcorecloud.ServiceClient, snapshotID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, snapshotID), &r.Body, nil)
	return
}

// ListAll is a convenience function that returns all snapshots.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Snapshot, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractSnapshots(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// IDFromName is a convenience function that returns a snapshot ID, given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string, opts ListOptsBuilder) (string, error) {
	count := 0
	id := ""

	all, err := ListAll(client, opts)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "snapshots"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "snapshots"}
	}
}
