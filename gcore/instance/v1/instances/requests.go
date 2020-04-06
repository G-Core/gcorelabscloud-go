package instances

import (
	"net/http"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/volume/v1/volumes"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToInstanceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	ExcludeSecGroup   string `q:"exclude_secgroup"`
	AvailableFloating bool   `q:"available_floating"`
}

// ToInstanceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToInstanceListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// DeleteOptsBuilder allows extensions to add additional parameters to the Delete request.
type DeleteOptsBuilder interface {
	ToInstanceDeleteQuery() (string, error)
}

// DeleteOpts. Set parameters for delete operation
type DeleteOpts struct {
	Volumes         []string `q:"volumes" validate:"omitempty,dive,uuid4" delimiter:"comma"`
	DeleteFloatings bool     `q:"delete_floatings" validate:"omitempty,allowed_without=FloatingIPs"`
	FloatingIPs     []string `q:"floatings" validate:"omitempty,allowed_without=DeleteFloatings,dive,uuid4" delimiter:"comma"`
}

// ToInstanceDeleteQuery formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToInstanceDeleteQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func (opts *DeleteOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// CreateVolumeOpts represents options used to create a volume.
type CreateVolumeOpts struct {
	Source     types.VolumeSource `json:"source" required:"true" validate:"required"`
	BootIndex  int                `json:"boot_index"`
	Size       int                `json:"size,omitempty" validate:"rfe=Source:image;new-volume"`
	TypeName   volumes.VolumeType `json:"type_name" required:"true" validate:"required"`
	Name       string             `json:"name,omitempty" validate:"omitempty"`
	ImageID    string             `json:"image_id,omitempty" validate:"rfe=Source:image,omitempty,uuid4"`
	SnapshotID string             `json:"snapshot_id,omitempty" validate:"rfe=Source:snapshot,omitempty,uuid4"`
	VolumeID   string             `json:"volume_id,omitempty" validate:"rfe=Source:existing-volume,omitempty,uuid4"`
}

func (opts *CreateVolumeOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

type CreateNewInterfaceFloatingIPOpts struct {
	Source             types.FloatingIPSource `json:"source" validate:"required,floating-ip-source"`
	ExistingFloatingID string                 `json:"existing_floating_id" validate:"rfe=Source:existing,omitempty,ip"`
}

// Validate
func (opts CreateNewInterfaceFloatingIPOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

type CreateInterfaceOpts struct {
	Type       types.InterfaceType               `json:"type" required:"true" validate:"required,interface-type"`
	NetworkID  string                            `json:"network_id,omitempty" validate:"rfe=Type:subnet,sfe=Type:external,omitempty,uuid4"`
	SubnetID   string                            `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,sfe=Type:external,omitempty,uuid4"`
	FloatingIP *CreateNewInterfaceFloatingIPOpts `json:"floating_ip,omitempty" validate:"omitempty,sfe=Type:external,dive"`
}

// Validate
func (opts CreateInterfaceOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// CreateOpts represents options used to create a instance.
type CreateOpts struct {
	Flavor         string                `json:"flavor" required:"true"`
	Names          []string              `json:"names,omitempty" validate:"required_without=NameTemplates"`
	NameTemplates  []string              `json:"name_templates,omitempty" validate:"required_without=Names"`
	Volumes        []CreateVolumeOpts    `json:"volumes" required:"true" validate:"required,dive"`
	Interfaces     []CreateInterfaceOpts `json:"interfaces" required:"true" validate:"required,dive"`
	SecurityGroups []gcorecloud.ItemID   `json:"security_groups" validate:"omitempty,dive,uuid4"`
	Keypair        string                `json:"keypair_name"`
	Password       string                `json:"password" validate:"omitempty,required_with=Username"`
	Username       string                `json:"username" validate:"omitempty,required_with=Password"`
	UserData       string                `json:"user_data" validate:"omitempty,base64"`
	Metadata       map[string]string     `json:"metadata,omitempty"`
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToInstanceCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// SecurityGroupOptsBuilder allows extensions to add parameters to the security groups request.
type SecurityGroupOptsBuilder interface {
	ToSecurityGroupActionMap() (map[string]interface{}, error)
}

type SecurityGroupOpts struct {
	Name string `json:"name" required:"true"`
}

// ToSecurityGroupActionMap builds a request body from SecurityGroupOpts.
func (opts SecurityGroupOpts) ToSecurityGroupActionMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
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
		return InstancePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific instance based on its unique ID.
func Get(client *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(client, id)
	var resp *http.Response
	resp, r.Err = client.Get(url, &r.Body, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}

// ListInterfaces retrieves network interfaces for instance
func ListInterfaces(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := interfacesListURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstanceInterfacePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all instances.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Instance, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstances(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// ListInterfacesAll is a convenience function that returns all instance interfaces.
func ListInterfacesAll(client *gcorecloud.ServiceClient, id string) ([]Interface, error) {
	pages, err := ListInterfaces(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstanceInterfaces(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// ListSecurityGroups retrieves security groups interfaces for instance
func ListSecurityGroups(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := securityGroupsListURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstanceSecurityGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListSecurityGroupsAll is a convenience function that returns all instance security groups.
func ListSecurityGroupsAll(client *gcorecloud.ServiceClient, id string) ([]gcorecloud.ItemIDName, error) {
	pages, err := ListSecurityGroups(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstanceSecurityGroups(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// AssignSecurityGroup adds a security groups to the instance.
func AssignSecurityGroup(client *gcorecloud.ServiceClient, id string, opts SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	var resp *http.Response
	resp, r.Err = client.Post(addSecurityGroupsURL(client, id), b, nil, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	defer func() {
		_ = resp.Body.Close()
	}()

	return
}

// UnAssignSecurityGroup removes a security groups from the instance.
func UnAssignSecurityGroup(client *gcorecloud.ServiceClient, id string, opts SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	var resp *http.Response
	resp, r.Err = client.Post(deleteSecurityGroupsURL(client, id), b, nil, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}

// Create creates an instance.
func Create(client *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var resp *http.Response
	resp, r.Err = client.Post(createURL(client, "v2"), b, &r.Body, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}

func Delete(client *gcorecloud.ServiceClient, instanceID string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := deleteURL(client, instanceID)
	if opts != nil {
		query, err := opts.ToInstanceDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	var resp *http.Response
	resp, r.Err = client.DeleteWithResponse(url, &r.Body, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}
