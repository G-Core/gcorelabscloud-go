package stacks

import (
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/stacks/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToStackCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// The name of the stack. It must start with an alphabetic character.
	Name string `json:"stack_name" required:"true"`
	// A structure that contains either the template file or url. Call the
	// associated methods to extract the information relevant to send in a create request.
	TemplateOpts *Template `json:"-" required:"true"`
	// Enables or disables deletion of all stack resources when a stack
	// creation fails. Default is true, meaning all resources are not deleted when
	// stack creation fails.
	DisableRollback *bool `json:"disable_rollback,omitempty"`
	// A structure that contains details for the environment of the stack.
	EnvironmentOpts *Environment `json:"-"`
	// User-defined parameters to pass to the template.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	// The timeout for stack creation in minutes.
	Timeout int `json:"timeout_mins,omitempty"`
	// A list of tags to assosciate with the Stack
	Tags []string `json:"-"`
}

// ToStackCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToStackCreateMap() (map[string]interface{}, error) {
	b, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if err := opts.TemplateOpts.Parse(); err != nil {
		return nil, err
	}

	if err := opts.TemplateOpts.getFileContents(opts.TemplateOpts.Parsed, ignoreIfTemplate, true); err != nil {
		return nil, err
	}
	opts.TemplateOpts.fixFileRefs()
	b["template"] = string(opts.TemplateOpts.Bin)

	files := make(map[string]string)
	for k, v := range opts.TemplateOpts.Files {
		files[k] = v
	}

	if opts.EnvironmentOpts != nil {
		if err := opts.EnvironmentOpts.Parse(); err != nil {
			return nil, err
		}
		if err := opts.EnvironmentOpts.getRRFileContents(ignoreIfEnvironment); err != nil {
			return nil, err
		}
		opts.EnvironmentOpts.fixFileRefs()
		for k, v := range opts.EnvironmentOpts.Files {
			files[k] = v
		}
		b["environment"] = string(opts.EnvironmentOpts.Bin)
	}

	if len(files) > 0 {
		b["files"] = files
	}

	if opts.Tags != nil {
		b["tags"] = strings.Join(opts.Tags, ",")
	}

	return b, nil
}

// Create accepts a CreateOpts struct and creates a new stack using the values
// provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToStackCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

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

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the Update operation in this package.
type UpdateOptsBuilder interface {
	ToStackUpdateMap() (map[string]interface{}, error)
}

// UpdatePatchOptsBuilder is the interface options structs have to satisfy in order
// to be used in the UpdatePatch operation in this package
type UpdatePatchOptsBuilder interface {
	ToStackUpdatePatchMap() (map[string]interface{}, error)
}

// UpdateOpts contains the common options struct used in this package's Update
// and UpdatePatch operations.
type UpdateOpts struct {
	// A structure that contains either the template file or url. Call the
	// associated methods to extract the information relevant to send in a create request.
	TemplateOpts *Template `json:"-"`
	// A structure that contains details for the environment of the stack.
	EnvironmentOpts *Environment `json:"-"`
	// User-defined parameters to pass to the template.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	// The timeout for stack creation in minutes.
	Timeout int `json:"timeout_mins,omitempty"`
	// A list of tags to associate with the Stack
	Tags []string `json:"-"`
}

// ToStackUpdateMap validates that a template was supplied and calls
// the toStackUpdateMap private function.
func (opts UpdateOpts) ToStackUpdateMap() (map[string]interface{}, error) {
	if opts.TemplateOpts == nil {
		return nil, ErrTemplateRequired{}
	}
	return toStackUpdateMap(opts)
}

// ToStackUpdatePatchMap calls the private function toStackUpdateMap
// directly.
func (opts UpdateOpts) ToStackUpdatePatchMap() (map[string]interface{}, error) {
	return toStackUpdateMap(opts)
}

// ToStackUpdateMap casts a CreateOpts struct to a map.
func toStackUpdateMap(opts UpdateOpts) (map[string]interface{}, error) {
	b, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	files := make(map[string]string)

	if opts.TemplateOpts != nil {
		if err := opts.TemplateOpts.Parse(); err != nil {
			return nil, err
		}

		if err := opts.TemplateOpts.getFileContents(opts.TemplateOpts.Parsed, ignoreIfTemplate, true); err != nil {
			return nil, err
		}
		opts.TemplateOpts.fixFileRefs()
		b["template"] = string(opts.TemplateOpts.Bin)

		for k, v := range opts.TemplateOpts.Files {
			files[k] = v
		}
	}

	if opts.EnvironmentOpts != nil {
		if err := opts.EnvironmentOpts.Parse(); err != nil {
			return nil, err
		}
		if err := opts.EnvironmentOpts.getRRFileContents(ignoreIfEnvironment); err != nil {
			return nil, err
		}
		opts.EnvironmentOpts.fixFileRefs()
		for k, v := range opts.EnvironmentOpts.Files {
			files[k] = v
		}
		b["environment"] = string(opts.EnvironmentOpts.Bin)
	}

	if len(files) > 0 {
		b["files"] = files
	}

	if opts.Tags != nil {
		b["tags"] = strings.Join(opts.Tags, ",")
	}

	return b, nil
}

// Update accepts an UpdateOpts struct and updates an existing stack using the
//  http PUT verb with the values provided. opts.TemplateOpts is required.
func Update(c *gcorecloud.ServiceClient, stackID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToStackUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, stackID), b, nil, nil)
	return
}

// Update accepts an UpdateOpts struct and updates an existing stack using the
//  http PATCH verb with the values provided. opts.TemplateOpts is not required.
func UpdatePatch(c *gcorecloud.ServiceClient, stackID string, opts UpdatePatchOptsBuilder) (r UpdateResult) {
	b, err := opts.ToStackUpdatePatchMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, stackID), b, nil, nil)
	return
}

// Delete deletes a stack based on the stack name and stack ID.
func Delete(c *gcorecloud.ServiceClient, stackID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, stackID), nil)
	return
}
