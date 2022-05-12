package instances

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

const DefaultAvailabilityZone = "nova"

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a instance resource.
func (r commonResult) Extract() (*Instance, error) {
	var s Instance
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type ListMetricsResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a instance's metrics resource.
func (r ListMetricsResult) Extract() ([]InstanceMetrics, error) {
	var s []InstanceMetrics
	err := r.ExtractInto(&s)
	return s, err
}

func (r ListMetricsResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "results")
}

type RemoteConsoleResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a remote console resource.
func (r RemoteConsoleResult) Extract() (*RemoteConsole, error) {
	var rc RemoteConsole
	err := r.ExtractInto(&rc)
	return &rc, err
}

func (r RemoteConsoleResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "remote_console")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Instance.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Instance.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Instance.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// MetadataActionResult represents the result of a create, delete or update operation(no content)
type MetadataActionResult struct {
	gcorecloud.ErrResult
}

// MetadataResult represents the result of a get operation
type MetadataResult struct {
	commonResult
}

// SecurityGroupActionResult represents the result of a actions operation(no content)
type SecurityGroupActionResult struct {
	gcorecloud.ErrResult
}

//  InterfaceActionResult represents the result of a actions operation(no content)
type InterfaceActionResult struct {
	gcorecloud.ErrResult
}

type RemoteConsole struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

type InstanceVolume struct {
	ID                  string `json:"id"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

type InstanceAddress struct {
	Type       types.AddressType `json:"type"`
	Address    net.IP            `json:"addr"`
	SubnetID   *string           `json:"subnet_id,omitempty"`
	SubnetName *string           `json:"subnet_name,omitempty"`
}

// InstanceMetrics represent a instance metrics struct
type InstanceMetrics struct {
	CPUUtil           float64                  `json:"cpu_util"`
	Disks             []DiskMetrics            `json:"disks"`
	MemoryUtil        float64                  `json:"memory_util"`
	NetworkBPSEgress  float64                  `json:"network_Bps_egress"`
	NetworkBPSIngress float64                  `json:"network_Bps_ingress"`
	NetworkPPSEgress  float64                  `json:"network_pps_egress"`
	NetworkPPSIngress float64                  `json:"network_pps_ingress"`
	Time              gcorecloud.JSONRFC3339ZZ `json:"time"`
}

type DiskMetrics struct {
	BpsRead   float64 `json:"disk_Bps_read"`
	BpsWrite  float64 `json:"disk_Bps_write"`
	IOPSRead  float64 `json:"disk_iops_read"`
	IOPSWrite float64 `json:"disk_iops_write"`
	Name      string  `json:"disk_name"`
}

// Instance represents a instance structure.
type Instance struct {
	ID               string                       `json:"instance_id"`
	Name             string                       `json:"instance_name"`
	Description      string                       `json:"instance_description"`
	CreatedAt        gcorecloud.JSONRFC3339ZZ     `json:"instance_created"`
	Status           string                       `json:"status"`
	VMState          string                       `json:"vm_state"`
	TaskState        *string                      `json:"task_state"`
	Flavor           flavors.Flavor               `json:"flavor"`
	Metadata         map[string]interface{}       `json:"metadata"`
	Volumes          []InstanceVolume             `json:"volumes"`
	Addresses        map[string][]InstanceAddress `json:"addresses"`
	SecurityGroups   []gcorecloud.ItemName        `json:"security_groups"`
	CreatorTaskID    *string                      `json:"creator_task_id"`
	TaskID           *string                      `json:"task_id"`
	ProjectID        int                          `json:"project_id"`
	RegionID         int                          `json:"region_id"`
	Region           string                       `json:"region"`
	AvailabilityZone string                       `json:"availability_zone"`
}

// Interface represents a instance port interface.
type Interface struct {
	PortID              string         `json:"port_id"`
	MacAddress          gcorecloud.MAC `json:"mac_address"`
	NetworkID           string         `json:"network_id"`
	PortSecurityEnabled bool           `json:"port_security_enabled"`
	IPAssignments       []PortIP       `json:"ip_assignments"`
	NetworkDetails      NetworkDetail  `json:"network_details"`
	FloatingIPDetails   []FloatingIP   `json:"floatingip_details"`
	SubPorts            []SubPort      `json:"sub_ports"`
}

// SubPort represent a instance sub port interface
type SubPort struct {
	PortID            string         `json:"port_id"`
	MacAddress        gcorecloud.MAC `json:"mac_address"`
	NetworkID         string         `json:"network_id"`
	IPAssignments     []PortIP       `json:"ip_assignments"`
	NetworkDetails    NetworkDetail  `json:"network_details"`
	FloatingIPDetails []FloatingIP   `json:"floatingip_details"`
	SegmentationID    int            `json:"segmentation_id"`
	SegmentationType  string         `json:"segmentation_type"`
}

// PortIP
type PortIP struct {
	IPAddress net.IP `json:"ip_address"`
	SubnetID  string `json:"subnet_id"`
}

// FloatingIP represents a floating ip of instance port.
type FloatingIP struct {
	FloatingIPAddress net.IP                   `json:"floating_ip_address"`
	RouterID          string                   `json:"router_id"`
	SubnetID          string                   `json:"subnet_id"`
	Status            string                   `json:"status"`
	ID                string                   `json:"id"`
	PortID            string                   `json:"port_id"`
	DNSDomain         string                   `json:"dns_domain"`
	DNSName           string                   `json:"dns_name"`
	FixedIPAddress    net.IP                   `json:"fixed_ip_address"`
	UpdatedAt         *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	CreatedAt         gcorecloud.JSONRFC3339Z  `json:"created_at"`
	CreatorTaskID     *string                  `json:"creator_task_id"`
	ProjectID         int                      `json:"project_id"`
	RegionID          int                      `json:"region_id"`
	Region            string                   `json:"region"`
}

// Subnet port subnet
type Subnet struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	IPVersion     gcorecloud.IPVersion     `json:"ip_version"`
	EnableDHCP    bool                     `json:"enable_dhcp"`
	Cidr          gcorecloud.CIDR          `json:"cidr"`
	CreatedAt     gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt     *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	NetworkID     string                   `json:"network_id"`
	TaskID        *string                  `json:"task_id"`
	CreatorTaskID *string                  `json:"creator_task_id"`
	ProjectID     int                      `json:"project_id"`
	RegionID      int                      `json:"region_id"`
	Region        string                   `json:"region"`
}

// NetworkDetail represents a NetworkDetails of instance port.
type NetworkDetail struct {
	Mtu           int                      `json:"mtu"`
	UpdatedAt     *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	CreatedAt     gcorecloud.JSONRFC3339Z  `json:"created_at"`
	ID            string                   `json:"id"`
	External      bool                     `json:"external"`
	Default       bool                     `json:"default"`
	Name          string                   `json:"name"`
	Shared        bool                     `json:"shared"`
	Subnets       []Subnet                 `json:"subnets"`
	ProjectID     int                      `json:"project_id"`
	RegionID      int                      `json:"region_id"`
	Region        string                   `json:"region"`
	TaskID        *string                  `json:"task_id"`
	CreatorTaskID *string                  `json:"creator_task_id"`
}

type InstancePorts struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	SecurityGroups []gcorecloud.ItemIDName `json:"security_groups"`
}

// InstancePage is the page returned by a pager when traversing over a
// collection of instances.
type InstancePage struct {
	pagination.LinkedPageBase
}

// MetadataPage is the page returned by a pager when traversing over a
// collection of instance metadata objects.
type MetadataPage struct {
	pagination.LinkedPageBase
}

// InstanceInterfacePage is the page returned by a pager when traversing over a
// collection of instance interfaces.
type InstanceInterfacePage struct {
	pagination.LinkedPageBase
}

// InstanceSecurityGroupPage is the page returned by a pager when traversing over a
// collection of instance security groups.
type InstanceSecurityGroupPage struct {
	pagination.LinkedPageBase
}

// InstancePortsPage is the page returned by a pager when traversing over a
// collection of instance ports.
type InstancePortsPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of instances has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r InstancePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// NextPageURL is invoked when a paginated collection of instance interfaces has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r InstanceInterfacePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// NextPageURL is invoked when a paginated collection of instance security groups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r InstanceSecurityGroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// NextPageURL is invoked when a paginated collection of instance ports has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r InstancePortsPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// NextPageURL is invoked when a paginated collection of instance metadata objects has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r MetadataPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a InstancePage struct is empty.
func (r InstancePage) IsEmpty() (bool, error) {
	is, err := ExtractInstances(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a InstanceInterfacePage struct is empty.
func (r InstanceInterfacePage) IsEmpty() (bool, error) {
	is, err := ExtractInstanceInterfaces(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a InstanceSecurityGroupPage struct is empty.
func (r InstanceSecurityGroupPage) IsEmpty() (bool, error) {
	is, err := ExtractInstanceSecurityGroups(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a InstancePortsPage struct is empty.
func (r InstancePortsPage) IsEmpty() (bool, error) {
	is, err := ExtractInstancePorts(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a MetadataPage struct is empty.
func (r MetadataPage) IsEmpty() (bool, error) {
	is, err := ExtractMetadata(r)
	return len(is) == 0, err
}

// ExtractInstances accepts a Page struct, specifically a InstancePage struct,
// and extracts the elements into a slice of instance structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractInstances(r pagination.Page) ([]Instance, error) {
	var s []Instance
	err := ExtractInstancesInto(r, &s)
	return s, err
}

// ExtractInstanceInterfaces accepts a Page struct, specifically a InstanceInterfacePage struct,
// and extracts the elements into a slice of instance interface structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractInstanceInterfaces(r pagination.Page) ([]Interface, error) {
	var s []Interface
	err := ExtractInstanceInterfacesInto(r, &s)
	return s, err
}

// ExtractInstanceSecurityGroups accepts a Page struct, specifically a InstanceSecurityGroupPage struct,
// and extracts the elements into a slice of instance security group structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractInstanceSecurityGroups(r pagination.Page) ([]gcorecloud.ItemIDName, error) {
	var s []gcorecloud.ItemIDName
	err := ExtractInstanceSecurityGroupInto(r, &s)
	return s, err
}

// ExtractInstancePorts accepts a Page struct, specifically a InstancePortsPage struct,
// and extracts the elements into a slice of instance security group structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractInstancePorts(r pagination.Page) ([]InstancePorts, error) {
	var s []InstancePorts
	err := ExtractInstancePortInto(r, &s)
	return s, err
}

// ExtractMetadata accepts a Page struct, specifically a MetadataPage struct,
// and extracts the elements into a slice of instance metadata structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractMetadata(r pagination.Page) ([]Metadata, error) {
	var s []Metadata
	err := ExtractMetadataInto(r, &s)
	return s, err
}

func ExtractInstancesInto(r pagination.Page, v interface{}) error {
	return r.(InstancePage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractInstanceInterfacesInto(r pagination.Page, v interface{}) error {
	return r.(InstanceInterfacePage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractInstanceSecurityGroupInto(r pagination.Page, v interface{}) error {
	return r.(InstanceSecurityGroupPage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractInstancePortInto(r pagination.Page, v interface{}) error {
	return r.(InstancePortsPage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractMetadataInto(r pagination.Page, v interface{}) error {
	return r.(MetadataPage).Result.ExtractIntoSlicePtr(v, "results")
}

// UnmarshalJSON - implements Unmarshaler interface
func (i *Instance) UnmarshalJSON(data []byte) error {
	i.AvailabilityZone = DefaultAvailabilityZone
	type Alias Instance
	tmp := (*Alias)(i)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	return nil
}

type InstanceTaskResult struct {
	Instances   []string `json:"instances"`
	Volumes     []string `json:"volumes"`
	FloatingIPs []string `json:"floating_ips"`
	Ports       []string `json:"ports"`
}

func ExtractInstanceIDFromTask(task *tasks.Task) (string, error) {
	var result InstanceTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode instance information in task structure: %w", err)
	}
	if len(result.Instances) == 0 {
		return "", fmt.Errorf("cannot decode instance information in task structure: %w", err)
	}
	return result.Instances[0], nil
}

func ExtractInstancePortIDFromTask(task *tasks.Task) (string, error) {
	var result InstanceTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode instance information in task structure: %w", err)
	}
	if len(result.Ports) == 0 {
		return "", fmt.Errorf("cannot decode instance information in task structure: %w", err)
	}
	return result.Ports[0], nil
}

type Metadata struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

// Extract is a function that accepts a result and extracts a instance metadata resource.
func (r MetadataResult) Extract() (*Metadata, error) {
	var s Metadata
	err := r.ExtractInto(&s)
	return &s, err
}

type InstanceLocation struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Name        string `json:"name"`
	RegionID    string `json:"region_id"`
	RegionName  string `json:"region_name"`
	ClientID    int    `json:"client_id"`
}

type SearchLocationResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a instance's location resource.
func (r SearchLocationResult) Extract() ([]InstanceLocation, error) {
	var s []InstanceLocation
	err := r.ExtractInto(&s)
	return s, err
}

func (r SearchLocationResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "results")
}
