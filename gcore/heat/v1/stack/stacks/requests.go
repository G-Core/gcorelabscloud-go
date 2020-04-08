package stacks

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/stacks/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToStackListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the network attributes you want to see returned.
type ListOpts struct {
	// TenantID is the UUID of the tenant. A tenant is also known as
	// a project.
	TenantID string `q:"tenant_id"`

	// ID filters the stack list by a stack ID
	ID string `q:"id"`

	// Status filters the stack list by a status.
	Status string `q:"status"`

	// Name filters the stack list by a name.
	Name string `q:"name"`

	// Marker is the ID of last-seen item.
	Marker string `q:"marker"`

	// Limit is an integer value for the limit of values to return.
	Limit int `q:"limit"`

	// SortKey allows you to sort by stack_name, stack_status, creation_time, or
	// update_time key.
	SortKey types.SortKey `q:"sort_keys"`

	// SortDir sets the direction, and is either `asc` or `desc`.
	SortDir types.SortDir `q:"sort_dir"`

	// AllTenants is a bool to show all tenants.
	AllTenants bool `q:"global_tenant"`

	// ShowDeleted set to `true` to include deleted stacks in the list.
	ShowDeleted bool `q:"show_deleted"`

	// ShowNested set to `true` to include nested stacks in the list.
	ShowNested bool `q:"show_nested"`

	// Tags lists stacks that contain one or more simple string tags.
	Tags string `q:"tags"`

	// TagsAny lists stacks that contain one or more simple string tags.
	TagsAny string `q:"tags_any"`

	// NotTags lists stacks that do not contain one or more simple string tags.
	NotTags string `q:"not_tags"`

	// NotTagsAny lists stacks that do not contain one or more simple string tags.
	NotTagsAny string `q:"not_tags_any"`
}

// ToStackListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToStackListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List all stacks
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToStackListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return StackPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// List all stacks
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]StackList, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractStacks(page)
}

// Get retrieves a specific heat stack.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
