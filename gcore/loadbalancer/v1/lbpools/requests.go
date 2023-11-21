package lbpools

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
		query, err := opts.ToLBPoolListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PoolPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific lbpool based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToLBPoolListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	LoadBalancerID *string `q:"loadbalancer_id"`
	ListenerID     *string `q:"listener_id"`
	MemberDetails  *bool   `q:"details"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLBPoolListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	ToLBPoolCreateMap() (map[string]interface{}, error)
}

// CreateMemberOptsBuilder allows extensions to add parameters to the CreateMember request.
type CreateMemberOptsBuilder interface {
	ToLBPoolMemberCreateMap() (map[string]interface{}, error)
}

// CreateSessionPersistenceOpts represents options used to create a lbpool listener pool session persistence rules.
type CreateSessionPersistenceOpts struct {
	PersistenceGranularity string                `json:"persistence_granularity,omitempty"`
	PersistenceTimeout     int                   `json:"persistence_timeout,omitempty"`
	Type                   types.PersistenceType `json:"type" required:"true"`
	CookieName             string                `json:"cookie_name,omitempty"`
}

// CreateHealthMonitorOptsBuilder allows extensions to add parameters to the CreateHealthMonitor request.
type CreateHealthMonitorOptsBuilder interface {
	ToHealthMonitorCreateMap() (map[string]interface{}, error)
}

// CreateHealthMonitorOpts represents options used to create a lbpool health monitor.
type CreateHealthMonitorOpts struct {
	ID             string                  `json:"id,omitempty"`
	Type           types.HealthMonitorType `json:"type" required:"true"`
	Delay          int                     `json:"delay" required:"true"`
	MaxRetries     int                     `json:"max_retries" required:"true"`
	Timeout        int                     `json:"timeout" required:"true"`
	MaxRetriesDown int                     `json:"max_retries_down,omitempty"`
	HTTPMethod     *types.HTTPMethod       `json:"http_method,omitempty"`
	URLPath        string                  `json:"url_path,omitempty"`
	ExpectedCodes  string                  `json:"expected_codes,omitempty"`
}

// CreatePoolMemberOpts represents options used to create a lbpool listener pool member.
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

// CreateOpts represents options used to create a lbpool.
type CreateOpts struct {
	Name                 string                        `json:"name" required:"true" validate:"required,name"`
	Protocol             types.ProtocolType            `json:"protocol" required:"true"`
	LBPoolAlgorithm      types.LoadBalancerAlgorithm   `json:"lb_algorithm" required:"true"`
	Members              []CreatePoolMemberOpts        `json:"members,omitempty"`
	LoadBalancerID       string                        `json:"loadbalancer_id,omitempty"`
	ListenerID           string                        `json:"listener_id,omitempty"`
	HealthMonitor        *CreateHealthMonitorOpts      `json:"healthmonitor,omitempty"`
	SessionPersistence   *CreateSessionPersistenceOpts `json:"session_persistence,omitempty"`
	TimeoutClientData    int                           `json:"timeout_client_data"`
	TimeoutMemberData    int                           `json:"timeout_member_data"`
	TimeoutMemberConnect int                           `json:"timeout_member_connect"`
}

// ToLBPoolCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToLBPoolCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToHealthMonitorCreateMap builds a request body from CreateHealthMonitorOpts.
func (opts CreateHealthMonitorOpts) ToHealthMonitorCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToLBPoolMemberCreateMap builds a request body from CreatePoolMemberOpts.
func (opts CreatePoolMemberOpts) ToLBPoolMemberCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new lbpool using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToLBPoolCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToLBPoolUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a lbpool.
type UpdateOpts struct {
	Name                 string                        `json:"name,omitempty"`
	Members              []CreatePoolMemberOpts        `json:"members,omitempty"`
	Protocol             types.ProtocolType            `json:"protocol,omitempty"`
	LBPoolAlgorithm      types.LoadBalancerAlgorithm   `json:"lb_algorithm,omitempty"`
	HealthMonitor        *CreateHealthMonitorOpts      `json:"healthmonitor,omitempty"`
	SessionPersistence   *CreateSessionPersistenceOpts `json:"session_persistence"`
	TimeoutClientData    int                           `json:"timeout_client_data"`
	TimeoutMemberData    int                           `json:"timeout_member_data"`
	TimeoutMemberConnect int                           `json:"timeout_member_connect"`
}

// ToLBPoolUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToLBPoolUpdateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing lbpool using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, lbpoolID string, opts UpdateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToLBPoolUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, lbpoolID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the lbpool associated with it.
func Delete(c *gcorecloud.ServiceClient, lbpoolID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, lbpoolID), &r.Body, nil)
	return
}

// ListAll returns all LB pools
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Pool, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractPools(page)
}

// CreateMember creates LB pool member
func CreateMember(c *gcorecloud.ServiceClient, lbpoolID string, opts CreateMemberOptsBuilder) (r tasks.Result) {
	b, err := opts.ToLBPoolMemberCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createMemberURL(c, lbpoolID), b, &r.Body, nil)
	return
}

// DeleteMember accepts a unique pool and member ID and deletes pool member.
func DeleteMember(c *gcorecloud.ServiceClient, lbpoolID string, memberID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteMemberURL(c, lbpoolID, memberID), &r.Body, nil)
	return
}

// DeleteHealthMonitor accepts a unique ID and deletes the lbpool's healthmonitor associated with it.
func DeleteHealthMonitor(c *gcorecloud.ServiceClient, lbpoolID string) (r DeleteHealthMonitorResult) {
	_, r.Err = c.Delete(healthMonitorURL(c, lbpoolID), &gcorecloud.RequestOpts{OkCodes: []int{http.StatusNoContent}})
	return
}

// CreateHealthMonitor creates LB pool healthmonitor
func CreateHealthMonitor(c *gcorecloud.ServiceClient, lbpoolID string, opts CreateHealthMonitorOptsBuilder) (r tasks.Result) {
	b, err := opts.ToHealthMonitorCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(healthMonitorURL(c, lbpoolID), b, &r.Body, nil)
	return
}
