package lbpools

import (
	"fmt"
	"net"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

type DeleteHealthMonitorResult struct {
	gcorecloud.ErrResult
}

// Extract is a function that accepts a result and extracts a pool resource.
func (r commonResult) Extract() (*Pool, error) {
	var s Pool
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractPoolMember is a function that accepts a result and extracts a pool member resource.
func (r commonResult) ExtractPoolMember() (*PoolMember, error) {
	var s PoolMember
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Pool.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Pool.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Pool.
type UpdateResult struct {
	commonResult
}

// PoolMember represents a pool member structure.
type PoolMember struct {
	Address            *net.IP                  `json:"address,omitempty"`
	ID                 string                   `json:"id"`
	Weight             int                      `json:"weight,omitempty"`
	SubnetID           string                   `json:"subnet_id,omitempty"`
	InstanceID         string                   `json:"instance_id,omitempty"`
	ProtocolPort       int                      `json:"protocol_port,omitempty"`
	ProvisioningStatus types.ProvisioningStatus `json:"provisioning_status,omitempty"`
	OperatingStatus    types.OperatingStatus    `json:"operating_status,omitempty"`
	MonitorAddress     net.IP                   `json:"monitor_address,omitempty"`
	MonitorPort        *int                     `json:"monitor_port,omitempty"`
}

// Pool represents a pool structure.
type Pool struct {
	LoadBalancers         []gcorecloud.ItemID         `json:"loadbalancers"`
	Listeners             []gcorecloud.ItemID         `json:"listeners"`
	SessionPersistence    *SessionPersistence         `json:"session_persistence"`
	LoadBalancerAlgorithm types.LoadBalancerAlgorithm `json:"lb_algorithm"`
	Name                  string                      `json:"name"`
	ID                    string                      `json:"id"`
	Protocol              types.ProtocolType          `json:"protocol"`
	Members               []PoolMember                `json:"members"`
	HealthMonitor         *HealthMonitor              `json:"healthmonitor"`
	ProvisioningStatus    types.ProvisioningStatus    `json:"provisioning_status"`
	OperatingStatus       types.OperatingStatus       `json:"operating_status"`
	CreatorTaskID         string                      `json:"creator_task_id"`
	TaskID                string                      `json:"task_id"`
	TimeoutClientData     int                         `json:"timeout_client_data"`
	TimeoutMemberData     int                         `json:"timeout_member_data"`
	TimeoutMemberConnect  int                         `json:"timeout_member_connect"`
}

// IsDeleted LB pool state.
func (p Pool) IsDeleted() bool {
	return p.ProvisioningStatus == types.ProvisioningStatusDeleted
}

// HealthMonitor for LB pool
type HealthMonitor struct {
	ID             string                  `json:"id"`
	Type           types.HealthMonitorType `json:"type"`
	Delay          int                     `json:"delay"`
	MaxRetries     int                     `json:"max_retries"`
	Timeout        int                     `json:"timeout"`
	MaxRetriesDown int                     `json:"max_retries_down,omitempty"`
	HTTPMethod     *types.HTTPMethod       `json:"http_method,omitempty"`
	URLPath        string                  `json:"url_path,omitempty"`
	ExpectedCodes  string                  `json:"expected_codes,omitempty"`
}

// SessionPersistenceOpts represents options used to create a lbpool listener pool session persistence rules.
type SessionPersistence struct {
	PersistenceGranularity string                `json:"persistence_granularity,omitempty"`
	PersistenceTimeout     int                   `json:"persistence_timeout,omitempty"`
	Type                   types.PersistenceType `json:"type"`
	CookieName             string                `json:"cookie_name,omitempty"`
}

// PoolPage is the page returned by a pager when traversing over a collection of pools.
type PoolPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of pools has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r PoolPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a PoolPage struct is empty.
func (r PoolPage) IsEmpty() (bool, error) {
	is, err := ExtractPools(r)
	return len(is) == 0, err
}

// ExtractPool accepts a Page struct, specifically a PoolPage struct,
// and extracts the elements into a slice of Pool structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPools(r pagination.Page) ([]Pool, error) {
	var s []Pool
	err := ExtractPoolsInto(r, &s)
	return s, err
}

func ExtractPoolsInto(r pagination.Page, v interface{}) error {
	return r.(PoolPage).Result.ExtractIntoSlicePtr(v, "results")
}

type PoolTaskResult struct {
	Pools []string `json:"pools"`
}

type HealthMonitorTaskResult struct {
	HealthMonitors []string `json:"healthmonitors"`
}

type PoolMemberTaskResult struct {
	Members []string `json:"members"`
}

func ExtractPoolIDFromTask(task *tasks.Task) (string, error) {
	var result PoolTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode pool information in task structure: %w", err)
	}
	if len(result.Pools) == 0 {
		return "", fmt.Errorf("cannot decode pool information in task structure: %w", err)
	}
	return result.Pools[0], nil
}

func ExtractHealthMonitorIDFromTask(task *tasks.Task) (string, error) {
	var result HealthMonitorTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode healthmonitor information in task structure: %w", err)
	}
	if len(result.HealthMonitors) == 0 {
		return "", fmt.Errorf("cannot decode healthmonitor information in task structure: %w", err)
	}
	return result.HealthMonitors[0], nil
}

func ExtractPoolMemberIDFromTask(task *tasks.Task) (string, error) {
	var result PoolMemberTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode pool member information in task structure: %w", err)
	}
	if len(result.Members) == 0 {
		return "", fmt.Errorf("cannot decode pool member information in task structure: %w", err)
	}
	return result.Members[0], nil
}
