package loadbalancers

import (
	"net"
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
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
	LoggingEnabled   bool              `q:"logging_enabled" validate:"omitempty"`
	MetadataK        string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV       map[string]string `q:"metadata_kv" validate:"omitempty"`
	WithDdos         bool              `q:"with_ddos" validate:"omitempty"`
	Name             string            `q:"name" validate:"omitempty"`
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
func Get(c *gcorecloud.ServiceClient, id string, opts GetOptsBuilder) (r GetResult) {
	url := getURL(c, id)
	if opts != nil {
		query, err := opts.ToLoadBalancerGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// GetOpts allows the filtering and sorting Get API response.
type GetOpts struct {
	ShowStats bool `q:"show_stats" validate:"omitempty"`
	WithDdos  bool `q:"with_ddos" validate:"omitempty"`
}

// ToLoadBalancerListQuery formats a ListOpts into a query string.
func (opts GetOpts) ToLoadBalancerGetQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}

	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToLoadBalancerGetQuery() (string, error)
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
	ExpectedCodes  string                  `json:"expected_codes,omitempty"`
}

// CreatePoolMemberOpts represents options used to create a loadbalancer listener pool member.
type CreatePoolMemberOpts struct {
	ID             string `json:"id,omitempty"`
	Address        net.IP `json:"address" required:"true"`
	ProtocolPort   int    `json:"protocol_port" required:"true"`
	Weight         int    `json:"weight,omitempty"`
	SubnetID       string `json:"subnet_id,omitempty"`
	InstanceID     string `json:"instance_id,omitempty"`
	MonitorAddress net.IP `json:"monitor_address,omitempty"`
	MonitorPort    *int   `json:"monitor_port,omitempty"`
}

// CreatePoolOpts represents options used to create a loadbalancer listener pool.
type CreatePoolOpts struct {
	Name                  string                        `json:"name" required:"true" validate:"required,name"`
	Protocol              types.ProtocolType            `json:"protocol" required:"true"`
	Members               []CreatePoolMemberOpts        `json:"members"`
	HealthMonitor         *CreateHealthMonitorOpts      `json:"healthmonitor,omitempty"`
	LoadBalancerAlgorithm types.LoadBalancerAlgorithm   `json:"lb_algorithm,omitempty"`
	SessionPersistence    *CreateSessionPersistenceOpts `json:"session_persistence,omitempty"`
	TimeoutClientData     *int                          `json:"timeout_client_data,omitempty"`
	TimeoutMemberData     *int                          `json:"timeout_member_data,omitempty"`
	TimeoutMemberConnect  *int                          `json:"timeout_member_connect,omitempty"`
}

// CreateListenerOpts represents options used to create a loadbalancer listener.
type CreateListenerOpts struct {
	Name                 string             `json:"name" required:"true" validate:"required,name"`
	ProtocolPort         int                `json:"protocol_port" required:"true"`
	Protocol             types.ProtocolType `json:"protocol" required:"true"`
	Certificate          string             `json:"certificate,omitempty"`
	CertificateChain     string             `json:"certificate_chain,omitempty"`
	PrivateKey           string             `json:"private_key,omitempty"`
	Pools                []CreatePoolOpts   `json:"pools,omitempty" validate:"omitempty,dive"`
	SecretID             string             `json:"secret_id,omitempty"`
	SNISecretID          []string           `json:"sni_secret_id,omitempty"`
	InsertXForwarded     bool               `json:"insert_x_forwarded"`
	AllowedCIDRS         []string           `json:"allowed_cidrs,omitempty" validate:"omitempty,dive,cidr"`
	TimeoutClientData    *int               `json:"timeout_client_data,omitempty"`
	TimeoutMemberData    *int               `json:"timeout_member_data,omitempty"`
	TimeoutMemberConnect *int               `json:"timeout_member_connect,omitempty"`
	ConnectionLimit      *int               `json:"connection_limit,omitempty"`
}

// CreateLoggingOpts represents options used to configure logging for a loadbalancer.
type CreateLoggingOpts struct {
	Enabled             bool                       `json:"enabled"`
	TopicName           string                     `json:"topic_name,omitempty"`
	DestinationRegionID int                        `json:"destination_region_id,omitempty"`
	RetentionPolicy     *CreateRetentionPolicyOpts `json:"retention_policy,omitempty" validate:"omitempty,dive"`
}

// CreateRetentionPolicyOpts represents options used to configure logging topic retention policy for a loadbalancer.
type CreateRetentionPolicyOpts struct {
	Period int `json:"period,omitempty"`
}

// CreateOpts represents options used to create a loadbalancer.
type CreateOpts struct {
	Name         string                                      `json:"name" required:"true" validate:"required,name"`
	Listeners    []CreateListenerOpts                        `json:"listeners,omitempty" validate:"omitempty,dive"`
	VipNetworkID string                                      `json:"vip_network_id,omitempty" validate:"omitempty,allowed_without=VipPortID"`
	VipSubnetID  string                                      `json:"vip_subnet_id,omitempty"`
	VipPortID    string                                      `json:"vip_port_id,omitempty" validate:"omitempty,allowed_without=VipNetworkID"`
	VIPIPFamily  types.IPFamilyType                          `json:"vip_ip_family,omitempty" validate:"omitempty,enum"`
	Flavor       *string                                     `json:"flavor,omitempty"`
	Tags         []string                                    `json:"tag,omitempty"`
	Metadata     map[string]string                           `json:"metadata,omitempty"`
	FloatingIP   *instances.CreateNewInterfaceFloatingIPOpts `json:"floating_ip,omitempty"`
	Logging      *CreateLoggingOpts                          `json:"logging,omitempty"`
}

// ToLoadBalancerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToLoadBalancerCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new loadbalancer using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder, reqOpts *gcorecloud.RequestOpts) (r tasks.Result) {
	b, err := opts.ToLoadBalancerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, reqOpts)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToLoadBalancerUpdateMap() (map[string]interface{}, error)
}

// UpdateLoggingOpts represents options used to configure logging for a loadbalancer.
type UpdateLoggingOpts struct {
	Enabled             bool                       `json:"enabled"`
	TopicName           string                     `json:"topic_name,omitempty"`
	DestinationRegionID int                        `json:"destination_region_id,omitempty"`
	RetentionPolicy     *UpdateRetentionPolicyOpts `json:"retention_policy,omitempty" validate:"omitempty,dive"`
}

// UpdateRetentionPolicyOpts represents options used to configure logging topic retention policy for a loadbalancer.
type UpdateRetentionPolicyOpts struct {
	Period int `json:"period,omitempty"`
}

// UpdateOpts represents options used to update a loadbalancer.
type UpdateOpts struct {
	Name    string             `json:"name,omitempty"`
	Logging *UpdateLoggingOpts `json:"logging,omitempty"`
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
func Delete(c *gcorecloud.ServiceClient, loadbalancerID string, reqOpts *gcorecloud.RequestOpts) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, loadbalancerID), &r.Body, reqOpts)
	return
}

// ToFloatingIPListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLoadBalancerQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
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

// ResizeOptsBuilder allows extensions to add additional parameters to the Resize request.
type ResizeOptsBuilder interface {
	ToLoadBalancerResizeMap() (map[string]interface{}, error)
}

// ResizeOpts represents options used to resize a loadbalancer.
type ResizeOpts struct {
	Flavor string `json:"flavor" required:"true"`
}

// ToLoadBalancerResizeMap builds a request body from ResizeOpts.
func (opts ResizeOpts) ToLoadBalancerResizeMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Resize accepts a ResizeOpts struct and resizes a loadbalancer using the values provided.
func Resize(c *gcorecloud.ServiceClient, loadbalancerID string, opts ResizeOptsBuilder, reqOpts *gcorecloud.RequestOpts) (r tasks.Result) {
	b, err := opts.ToLoadBalancerResizeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(resizeLoadBalancerUrl(c, loadbalancerID), b, &r.Body, reqOpts)
	return
}
