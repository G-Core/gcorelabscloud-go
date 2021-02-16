package bminstances

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToInstanceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Name     string `q:"name"`
	FlavorID string `q:"flavor_id"`
}

// ToInstanceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToInstanceListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

type CreateNewInterfaceFloatingIPOpts struct {
	Source             types.FloatingIPSource `json:"source" validate:"required,enum"`
	ExistingFloatingID string                 `json:"existing_floating_id" validate:"rfe=Source:existing,sfe=Source:new,omitempty,uuid4"`
}

// Validate
func (opts CreateNewInterfaceFloatingIPOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

type InterfaceOpts struct {
	Type       types.InterfaceType               `json:"type" validate:"omitempty,enum"`
	NetworkID  string                            `json:"network_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	SubnetID   string                            `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	FloatingIP *CreateNewInterfaceFloatingIPOpts `json:"floating_ip,omitempty" validate:"omitempty,dive"`
}

// Validate
func (opts InterfaceOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a instance.
type CreateOpts struct {
	Flavor        string                 `json:"flavor" required:"true"`
	Names         []string               `json:"names,omitempty" validate:"required_without=NameTemplates"`
	NameTemplates []string               `json:"name_templates,omitempty" validate:"required_without=Names"`
	ImageID       string                 `json:"image_id,omitempty" validate:"required_without=AppTemplateID"`
	AppTemplateID string                 `json:"apptemplate_id,omitempty" validate:"required_without=ImageID"`
	Interfaces    []InterfaceOpts        `json:"interfaces" required:"true" validate:"required,dive"`
	Keypair       string                 `json:"keypair_name"`
	Password      string                 `json:"password" validate:"omitempty,required_with=Username"`
	Username      string                 `json:"username" validate:"omitempty,required_with=Password"`
	UserData      string                 `json:"user_data,omitempty" validate:"omitempty,base64"`
	AppConfig     map[string]interface{} `json:"app_config,omitempty" validate:"omitempty"`
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToInstanceCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	var err error
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if len(opts.AppConfig) > 0 {
		mp["app_config"] = opts.AppConfig
	} else {
		delete(mp, "app_config")
	}
	return mp, nil
}

// Create creates an baremetal instance.
func Create(client *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, nil) // nolint
	return
}

func List(client *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToInstanceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return instances.InstancePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all instances.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]instances.Instance, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := instances.ExtractInstances(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}
