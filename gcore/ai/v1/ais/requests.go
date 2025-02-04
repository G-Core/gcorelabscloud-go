package ai

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// DeleteOptsBuilder allows extensions to add additional parameters to the Delete request.
type DeleteOptsBuilder interface {
	ToAIClusterDeleteQuery() (string, error)
}

// DeleteOpts. Set parameters for delete operation
type DeleteOpts struct {
	Volumes          []string `q:"volumes" validate:"omitempty,dive,uuid4" delimiter:"comma"`
	DeleteFloatings  bool     `q:"delete_floatings" validate:"omitempty,allowed_without=FloatingIPs"`
	FloatingIPs      []string `q:"floatings" validate:"omitempty,allowed_without=DeleteFloatings,dive,uuid4" delimiter:"comma"`
	ReservedFixedIPs []string `q:"reserved_fixed_ips" validate:"omitempty,dive,uuid4" delimiter:"comma"`
}

// ToIAClusterDeleteQuery formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToAIClusterDeleteQuery() (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func (opts *DeleteOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToAIClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a AI Cluster.
type CreateOpts struct {
	Flavor         string                                  `json:"flavor" validate:"required,min=1"`
	Name           string                                  `json:"name" validate:"required,min=3,max=63"`
	ImageID        string                                  `json:"image_id" validate:"required,uuid4"`
	Interfaces     []instances.InterfaceInstanceCreateOpts `json:"interfaces" validate:"required,dive"`
	Volumes        []instances.CreateVolumeOpts            `json:"volumes,omitempty" validate:"omitempty,required,dive"`
	SecurityGroups []gcorecloud.ItemID                     `json:"security_groups,omitempty" validate:"omitempty,dive,uuid4"`
	Keypair        string                                  `json:"keypair_name,omitempty"`
	Password       string                                  `json:"password" validate:"omitempty,required_with=Username"`
	Username       string                                  `json:"username" validate:"omitempty,required_with=Password"`
	UserData       string                                  `json:"user_data,omitempty" validate:"omitempty,base64"`
	Metadata       map[string]string                       `json:"metadata,omitempty" validate:"omitempty,dive"`
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToAIClusterCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToAIClusterCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// ResizeAIClusterOptsBuilder builds parameters or change flavor request.
type ResizeAIClusterOptsBuilder interface {
	ToResizeAIClusterActionMap() (map[string]interface{}, error)
}

// ResizeAIClusterOpts represents options used to resize AI Clsuter.
type ResizeAIClusterOpts struct {
	Flavor         string                                  `json:"flavor" validate:"omitempty,min=1"`
	ImageID        string                                  `json:"image_id" validate:"omitempty,uuid4"`
	Interfaces     []instances.InterfaceInstanceCreateOpts `json:"interfaces" validate:"required,dive"`
	Volumes        []instances.CreateVolumeOpts            `json:"volumes,omitempty" validate:"omitempty,dive"`
	SecurityGroups []gcorecloud.ItemID                     `json:"security_groups,omitempty" validate:"omitempty,dive,uuid4"`
	Keypair        string                                  `json:"keypair_name,omitempty"`
	Password       string                                  `json:"password" validate:"omitempty,required_with=Username"`
	Username       string                                  `json:"username" validate:"omitempty,required_with=Password"`
	UserData       string                                  `json:"user_data,omitempty" validate:"omitempty,base64"`
	Metadata       map[string]string                       `json:"metadata,omitempty" validate:"omitempty,dive"`
}

// Validate
func (opts ResizeAIClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToResizeAIClusterActionMap builds a request body from ResizeAIClusterOpts.
func (opts ResizeAIClusterOpts) ToResizeAIClusterActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// AttachInterfaceOptsBuilder allows extensions to add parameters to the interface request.
type AttachInterfaceOptsBuilder interface {
	ToInterfaceAttachMap() (map[string]interface{}, error)
}

type AttachInterfaceOpts struct {
	Type      types.InterfaceType `json:"type,omitempty" validate:"omitempty,enum"`
	NetworkID string              `json:"network_id,omitempty" validate:"rfe=Type:any_subnet,omitempty,uuid4"`
	SubnetID  string              `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	PortID    string              `json:"port_id,omitempty" validate:"rfe=Type:reserved_fixed_ip,allowed_without_all=NetworkID SubnetID,omitempty,uuid4"`
	IpAddress string              `json:"ip_address,omitempty" validate:"allowed_without_all=Type NetworkID SubnetID FloatingIP,omitempty"`
}

// Validate
func (opts AttachInterfaceOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToInterfaceAttachMap builds a request body from AttachInterfaceOpts.
func (opts AttachInterfaceOpts) ToInterfaceAttachMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// DetachInterfaceOptsBuilder allows extensions to add parameters to the interface request.
type DetachInterfaceOptsBuilder interface {
	ToInterfaceDetachMap() (map[string]interface{}, error)
}

type DetachInterfaceOpts struct {
	PortID    string `json:"port_id" validate:"required,uuid4"`
	IpAddress string `json:"ip_address" validate:"required,ip"`
}

// Validate
func (opts DetachInterfaceOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToInterfaceDetachMap builds a request body from InterfaceOpts.
func (opts DetachInterfaceOpts) ToInterfaceDetachMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

func List(client *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AIClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all AI Clusters
func ListAll(client *gcorecloud.ServiceClient) ([]AICluster, error) {
	pages, err := List(client).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractAIClusters(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Get retrieves a specific AI Cluster based on its unique ID.
func Get(client *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(client, id)
	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}

// ListInterfaces retrieves network interfaces for AI Cluster
func ListInterfaces(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := interfacesListURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AIClusterInterfacePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListInterfacesAll is a convenience function that returns all AI Cluster interfaces.
func ListInterfacesAll(client *gcorecloud.ServiceClient, id string) ([]Interface, error) {
	pages, err := ListInterfaces(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractAIClusterInterfaces(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// ListPorts retrieves ports for AI Cluster
func ListPorts(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := portsListURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AIClusterPortsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListPortsAll is a convenience function that returns all AI Cluster ports.
func ListPortsAll(client *gcorecloud.ServiceClient, id string) ([]AIClusterPort, error) {
	pages, err := ListPorts(client, id).AllPages()
	if err != nil {
		return nil, err
	}
	all, err := ExtractAIClusterPorts(pages)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// AssignSecurityGroup adds a security groups to the AI Cluster
func AssignSecurityGroup(client *gcorecloud.ServiceClient, id string, opts instances.SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(addSecurityGroupsURL(client, id), b, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// UnAssignSecurityGroup removes a security groups from the AI Cluster
func UnAssignSecurityGroup(client *gcorecloud.ServiceClient, id string, opts instances.SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(deleteSecurityGroupsURL(client, id), b, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// AttachInterface adds a interface to the AI instance
func AttachAIInstanceInterface(client *gcorecloud.ServiceClient, instance_id string, opts AttachInterfaceOptsBuilder) (r tasks.Result) {
	b, err := opts.ToInterfaceAttachMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(attachAIInstanceInterfaceURL(client, instance_id), b, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// DetachInterface removes a interface from the AI instance
func DetachAIInstanceInterface(client *gcorecloud.ServiceClient, instance_id string, opts DetachInterfaceOptsBuilder) (r tasks.Result) {
	b, err := opts.ToInterfaceDetachMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(detachAIInstanceInterfaceURL(client, instance_id), b, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// Create creates an AI Cluster
func Create(client *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToAIClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, nil) // nolint
	return
}

// Delete an AI Cluster
func Delete(client *gcorecloud.ServiceClient, instanceID string, opts DeleteOptsBuilder) (r tasks.Result) {
	url := deleteURL(client, instanceID)
	if opts != nil {
		query, err := opts.ToAIClusterDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
	return
}

// PowerCycle AI instance.
func PowerCycleAIInstance(client *gcorecloud.ServiceClient, instance_id string) (r AIInstanceActionResult) {
	_, r.Err = client.Post(powerCycleAIInstanceURL(client, instance_id), nil, &r.Body, nil) // nolint
	return
}

// PowerCycle AI Cluster.
func PowerCycleAICluster(client *gcorecloud.ServiceClient, id string) (r AIClusterActionResult) {
	_, r.Err = client.Post(powerCycleAIURL(client, id), nil, &r.Body, nil) // nolint
	return
}

// Reboot AI instance.
func RebootAIInstance(client *gcorecloud.ServiceClient, instance_id string) (r AIInstanceActionResult) {
	_, r.Err = client.Post(rebootAIInstanceURL(client, instance_id), nil, &r.Body, nil) // nolint
	return
}

// Reboot AI cluster.
func RebootAICluster(client *gcorecloud.ServiceClient, id string) (r AIClusterActionResult) {
	_, r.Err = client.Post(rebootAIURL(client, id), nil, &r.Body, nil) // nolint
	return
}

// Suspend AI Cluster.
func Suspend(client *gcorecloud.ServiceClient, id string) (r tasks.Result) {
	_, r.Err = client.Post(suspendAIURL(client, id), nil, &r.Body, nil) // nolint
	return
}

// Resume AI Cluster.
func Resume(client *gcorecloud.ServiceClient, id string) (r tasks.Result) {
	_, r.Err = client.Post(resumeAIURL(client, id), nil, &r.Body, nil) // nolint
	return
}

// Resize AI Cluster.
func Resize(client *gcorecloud.ServiceClient, id string, opts ResizeAIClusterOptsBuilder) (r tasks.Result) {
	b, err := opts.ToResizeAIClusterActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(resizeAIURL(client, id), b, &r.Body, nil) // nolint
	return
}

func MetadataList(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := metadataURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return MetadataPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func MetadataListAll(client *gcorecloud.ServiceClient, id string) ([]metadata.Metadata, error) {
	pages, err := MetadataList(client, id).AllPages()
	if err != nil {
		return nil, err
	}
	all, err := ExtractMetadata(pages)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// MetadataCreateOrUpdate creates or update a metadata for an AI.
func MetadataCreateOrUpdate(client *gcorecloud.ServiceClient, id string, opts map[string]interface{}) (r MetadataActionResult) {
	_, r.Err = client.Post(metadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataReplace replace a metadata for an AI Cluster.
func MetadataReplace(client *gcorecloud.ServiceClient, id string, opts map[string]interface{}) (r MetadataActionResult) {
	_, r.Err = client.Put(metadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataDelete deletes defined metadata key for a AI Cluster.
func MetadataDelete(client *gcorecloud.ServiceClient, id string, key string) (r MetadataActionResult) {
	_, r.Err = client.Delete(metadataItemURL(client, id, key), &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataGet gets defined metadata key for a AI Cluster.
func MetadataGet(client *gcorecloud.ServiceClient, id string, key string) (r MetadataResult) {
	url := metadataItemURL(client, id, key)

	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}

// GetInstanceConsole retrieves a specific spice console based on instance unique ID.
func GetInstanceConsole(client *gcorecloud.ServiceClient, id string) (r RemoteConsoleResult) {
	url := getAIInstanceConsoleURL(client, id)
	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}

// RebuildGPUAIClusterOptsBuilder allows extensions to add additional parameters to the Rebuild request.
type RebuildGPUAIClusterOptsBuilder interface {
	ToRebuildGPUAIClusterMap() (map[string]interface{}, error)
}

// RebuildGPUAIClusterOpts represents options used to rebuild a GPU AI cluster.
type RebuildGPUAIClusterOpts struct {
	ImageID  string   `json:"image_id,omitempty"`
	UserData string   `json:"user_data,omitempty"`
	Nodes    []string `json:"nodes" required:"true" validate:"required"`
}

// Validate implements basic validation for this request.
func (opts RebuildGPUAIClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToRebuildGPUAIClusterMap builds a request body from RebuildGPUAIClusterOpts.
func (opts RebuildGPUAIClusterOpts) ToRebuildGPUAIClusterMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// RebuildGPUAICluster rebuilds a GPU AI cluster.
func RebuildGPUAICluster(client *gcorecloud.ServiceClient, clusterID string, opts RebuildGPUAIClusterOptsBuilder) (r tasks.Result) {
	b, err := opts.ToRebuildGPUAIClusterMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rebuildGPUAIURL(client, clusterID), b, &r.Body, nil) // nolint
	return
}
