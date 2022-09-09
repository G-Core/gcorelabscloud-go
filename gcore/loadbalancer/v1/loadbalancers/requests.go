package loadbalancers

import (
	"net"
	"net/http"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToLoadBalancerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return LoadBalancerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListOpts allows the filtering and sorting List API response.
type ListOpts struct {
	ShowStats        bool              `q:"show_stats" validate:"omitempty"`
	AssignedFloating bool              `q:"assigned_floating" validate:"omitempty"`
	MetadataK        string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV       map[string]string `q:"metadata_kv" validate:"omitempty"`
}

// ToLoadBalancerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLoadBalancerListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}

	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToLoadBalancerListQuery() (string, error)
}

// Get retrieves a specific loadbalancer based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToLoadBalancerCreateMap() (map[string]interface{}, error)
}

// CreateSessionPersistenceOpts represents options used to create a loadbalancer listener pool session persistence rules.
type CreateSessionPersistenceOpts struct {
	PersistenceGranularity string                `json:"persistence_granularity,omitempty"`
	PersistenceTimeout     int                   `json:"persistence_timeout,omitempty"`
	Type                   types.PersistenceType `json:"type" required:"true"`
	CookieName             string                `json:"cookie_name,omitempty"`
}

// CreateHealthMonitorOpts represents options used to create a loadbalancer health monitor.
type CreateHealthMonitorOpts struct {
	Type           types.HealthMonitorType `json:"type" required:"true"`
	Delay          int                     `json:"delay" required:"true"`
	MaxRetries     int                     `json:"max_retries" required:"true"`
	Timeout        int                     `json:"timeout" required:"true"`
	MaxRetriesDown int                     `json:"max_retries_down,omitempty"`
	HTTPMethod     *types.HTTPMethod       `json:"http_method,omitempty"`
	URLPath        string                  `json:"url_path,omitempty"`
}

// CreatePoolMemberOpts represents options used to create a loadbalancer listener pool member.
type CreatePoolMemberOpts struct {
	ID           string `json:"id,omitempty"`
	Address      net.IP `json:"address" required:"true"`
	ProtocolPort int    `json:"protocol_port" required:"true"`
	Weight       int    `json:"weight,omitempty"`
	SubnetID     string `json:"subnet_id,omitempty"`
	InstanceID   string `json:"instance_id,omitempty"`
}

// CreatePoolOpts represents options used to create a loadbalancer listener pool.
type CreatePoolOpts struct {
	Name                  string                        `json:"name" required:"true" validate:"required,name"`
	Protocol              types.ProtocolType            `json:"protocol" required:"true"`
	Members               []CreatePoolMemberOpts        `json:"members"`
	HealthMonitor         *CreateHealthMonitorOpts      `json:"healthmonitor,omitempty"`
	LoadBalancerAlgorithm types.LoadBalancerAlgorithm   `json:"lb_algorithm,omitempty"`
	SessionPersistence    *CreateSessionPersistenceOpts `json:"session_persistence,omitempty"`
}

// CreateListenerOpts represents options used to create a loadbalancer listener.
type CreateListenerOpts struct {
	Name             string             `json:"name" required:"true" validate:"required,name"`
	ProtocolPort     int                `json:"protocol_port" required:"true"`
	Protocol         types.ProtocolType `json:"protocol" required:"true"`
	Certificate      string             `json:"certificate,omitempty"`
	CertificateChain string             `json:"certificate_chain,omitempty"`
	PrivateKey       string             `json:"private_key,omitempty"`
	Pools            []CreatePoolOpts   `json:"pools,omitempty" validate:"omitempty,dive"`
	SecretID         string             `json:"secret_id,omitempty"`
	SNISecretID      []string           `json:"sni_secret_id,omitempty"`
	InsertXForwarded bool               `json:"insert_x_forwarded"`
}

// CreateOpts represents options used to create a loadbalancer.
type CreateOpts struct {
	Name         string                 `json:"name" required:"true" validate:"required,name"`
	Listeners    []CreateListenerOpts   `json:"listeners,omitempty" validate:"omitempty,dive"`
	VipNetworkID string                 `json:"vip_network_id,omitempty"`
	VipSubnetID  string                 `json:"vip_subnet_id,omitempty"`
	Flavor       *string                `json:"flavor,omitempty"`
	Tags         []string               `json:"tag,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ToLoadBalancerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToLoadBalancerCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new loadbalancer using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToLoadBalancerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToLoadBalancerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a loadbalancer.
type UpdateOpts struct {
	Name string `json:"name,omitempty" required:"true" validate:"required,name"`
}

// ToLoadBalancerUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToLoadBalancerUpdateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing loadbalancer using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, loadbalancerID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToLoadBalancerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, loadbalancerID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})
	return
}

// Delete accepts a unique ID and deletes the loadbalancer associated with it.
func Delete(c *gcorecloud.ServiceClient, loadbalancerID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, loadbalancerID), &r.Body, nil)
	return
}

// ListAll returns all LBs
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]LoadBalancer, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractLoadBalancers(page)
}

// CreateCustomSecurityGroup accepts a unique ID and create a custom security group for the load balancer's ingress port.
func CreateCustomSecurityGroup(c *gcorecloud.ServiceClient, loadbalancerID string) (r CustomSecurityGroupCreateResult) {
	_, r.Err = c.Post(createCustomSecurityGroupURL(c, loadbalancerID), nil, nil, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	return
}

// ListCustomSecurityGroup accepts a unique ID and returns a custom security group for the load balancer's ingress port.
func ListCustomSecurityGroup(c *gcorecloud.ServiceClient, loadbalancerID string) (r CustomSecurityGroupGetResult) {
	_, r.Err = c.Get(createCustomSecurityGroupURL(c, loadbalancerID), &r.Body, nil)
	return
}
