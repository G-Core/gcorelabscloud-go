package ai

import (
	"encoding/json"
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a AI Cluster resource.
func (r commonResult) Extract() (*AICluster, error) {
	var s AICluster
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type RemoteConsole struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
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
// method to interpret it as a AI Cluster.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a AI Cluster.
type GetResult struct {
	commonResult
}

type aiInstanceResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a AI instance action resource.
func (r aiInstanceResult) Extract() (*instances.Instance, error) {
	var s instances.Instance
	err := r.ExtractInto(&s)
	return &s, err
}

func (r aiInstanceResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// AIInstanceActionResult represents the result of an cluster instance operation. Call its Extract
// method to interpret it as a AI Cluster instance actions.
type AIInstanceActionResult struct {
	aiInstanceResult
}

type aiClusterResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a AI Cluster action resource.
func (r aiClusterResult) Extract() ([]instances.Instance, error) {
	var s []instances.Instance
	err := r.ExtractInto(&s)
	return s, err
}

func (r aiClusterResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "results")
}

// AIClusterActionResult represents the result of an cluster operation. Call its Extract
// method to interpret it as a AI Cluster actions.
type AIClusterActionResult struct {
	aiClusterResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// MetadataActionResult represents the result of a create, delete or update operation
type MetadataActionResult struct {
	gcorecloud.ErrResult
}

// MetadataResult represents the result of a get operation
type MetadataResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a AI Clsuter metadata resource.
func (r MetadataResult) Extract() (*metadata.Metadata, error) {
	var s metadata.Metadata
	err := r.ExtractInto(&s)
	return &s, err
}

// SecurityGroupActionResult represents the result of a actions operation
type SecurityGroupActionResult struct {
	gcorecloud.ErrResult
}

type PoplarInterfaceSecGrop struct {
	PortID         string                `json:"port_id"`
	NetworkID      string                `json:"network_id"`
	SecurityGroups []gcorecloud.ItemName `json:"security_groups"`
}

type AIClusterInterface struct {
	Type      string `json:"type"`
	NetworkID string `json:"network_id"`
	SubnetID  string `json:"subnet_id"`
	PortID    string `json:"port_id"`
}

// AICluster represents a AI Cluster structure.
type AICluster struct {
	ClusterID      string                         `json:"cluster_id"`
	ClusterName    string                         `json:"cluster_name"`
	ClusterStatus  string                         `json:"cluster_status"`
	TaskID         *string                        `json:"task_id"`
	TaskStatus     string                         `json:"task_status"`
	CreatedAt      gcorecloud.JSONRFC3339MilliNoZ `json:"instance_created"`
	ImageID        string                         `json:"image_id"`
	ImageName      string                         `json:"image_name"`
	Flavor         string                         `json:"flavor"`
	Volumes        []volumes.Volume               `json:"volumes"`
	SecurityGroups []PoplarInterfaceSecGrop       `json:"security_groups"`
	Interfaces     []AIClusterInterface           `json:"interfaces"`
	KeypairName    string                         `json:"keypair_name"`
	UserData       string                         `json:"user_data"`
	Username       string                         `json:"username"`
	Password       string                         `json:"password"`
	PoplarServer   []instances.Instance           `json:"poplar_servers"`
	Metadata       map[string]interface{}         `json:"cluster_metadata"`
	ProjectID      int                            `json:"project_id"`
	RegionID       int                            `json:"region_id"`
	Region         string                         `json:"region"`
}

// Interface represents a AI Cluster interface.
type Interface struct {
	PortID              string                  `json:"port_id"`
	MacAddress          gcorecloud.MAC          `json:"mac_address"`
	NetworkID           string                  `json:"network_id"`
	PortSecurityEnabled bool                    `json:"port_security_enabled"`
	IPAssignments       []instances.PortIP      `json:"ip_assignments"`
	NetworkDetails      instances.NetworkDetail `json:"network_details"`
	FloatingIPDetails   []instances.FloatingIP  `json:"floatingip_details"`
	SubPorts            []instances.SubPort     `json:"sub_ports"`
}

type AIClusterPort struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	SecurityGroups []gcorecloud.ItemIDName `json:"security_groups"`
}

// AIClusterPage is the page returned by a pager when traversing over a
// collection of ai clusters.
type AIClusterPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of AI Clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r AIClusterPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a AIClusterPage struct is empty.
func (r AIClusterPage) IsEmpty() (bool, error) {
	is, err := ExtractAIClusters(r)
	return len(is) == 0, err
}

// AIClusterInterfacePage is the page returned by a pager when traversing over a
// collection of cluster interfaces.
type AIClusterInterfacePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of AI Cluster interfaces has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r AIClusterInterfacePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a AIClusterInterfacePage struct is empty.
func (r AIClusterInterfacePage) IsEmpty() (bool, error) {
	is, err := ExtractAIClusterInterfaces(r)
	return len(is) == 0, err
}

// AIClusterPortsPage is the page returned by a pager when traversing over a
// collection of ai cluster ports.
type AIClusterPortsPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of ai cluster ports has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r AIClusterPortsPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a AIClusterPortsPage struct is empty.
func (r AIClusterPortsPage) IsEmpty() (bool, error) {
	is, err := ExtractAIClusterPorts(r)
	return len(is) == 0, err
}

// MetadataPage is the page returned by a pager when traversing over a
// collection of AI Cluster metadata objects.
type MetadataPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of AI Cluster metadata objects has reached
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

// IsEmpty checks whether a MetadataPage struct is empty.
func (r MetadataPage) IsEmpty() (bool, error) {
	is, err := ExtractMetadata(r)
	return len(is) == 0, err
}

// ExtractAIClusters accepts a Page struct, specifically a AIClustersPage struct,
// and extracts the elements into a slice of AICluster structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAIClusters(r pagination.Page) ([]AICluster, error) {
	var s []AICluster
	err := ExtractAIClustersInto(r, &s)
	return s, err
}

func ExtractAIClustersInto(r pagination.Page, v interface{}) error {
	return r.(AIClusterPage).Result.ExtractIntoSlicePtr(v, "results")
}

// ExtractAIClusterInterfaces accepts a Page struct, specifically a AIClusterInterfacePage struct,
// and extracts the elements into a slice of AI Cluster interface structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAIClusterInterfaces(r pagination.Page) ([]Interface, error) {
	var s []Interface
	err := ExtractAIClusterInterfacesInto(r, &s)
	return s, err
}

// ExtractAIClusterInterfacesInto accepts a Page struct, specifically a AIClusterInterfacePage struct,
// and extracts the elements into a slice of AI Cluster interface structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAIClusterInterfacesInto(r pagination.Page, v interface{}) error {
	return r.(AIClusterInterfacePage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractAIClusterPortInto(r pagination.Page, v interface{}) error {
	return r.(AIClusterPortsPage).Result.ExtractIntoSlicePtr(v, "results")
}

// ExtractAIClusterPorts accepts a Page struct, specifically a AIClusterPortsPage struct,
// and extracts the elements into a slice of cluster porst structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAIClusterPorts(r pagination.Page) ([]AIClusterPort, error) {
	var s []AIClusterPort
	err := ExtractAIClusterPortInto(r, &s)
	return s, err
}

// ExtractMetadata accepts a Page struct, specifically a MetadataPage struct,
// and extracts the elements into a slice of ai cluster metadata structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractMetadata(r pagination.Page) ([]metadata.Metadata, error) {
	var s []metadata.Metadata
	err := ExtractMetadataInto(r, &s)
	return s, err
}

func ExtractMetadataInto(r pagination.Page, v interface{}) error {
	return r.(MetadataPage).Result.ExtractIntoSlicePtr(v, "results")
}

// UnmarshalJSON - implements Unmarshaler interface
func (i *AICluster) UnmarshalJSON(data []byte) error {
	type Alias AICluster
	tmp := (*Alias)(i)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	return nil
}

type AIClusterTaskResult struct {
	AIClusters []string `mapstructure:"ai_clusters"`
	// etc
}

func ExtractAIClusterIDFromTask(task *tasks.Task) (string, error) {
	var result AIClusterTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode AI cluster information in task structure: %w", err)
	}
	if len(result.AIClusters) == 0 {
		return "", fmt.Errorf("cannot decode ai cluster information in task structure: %w", err)
	}
	return result.AIClusters[0], nil
}
