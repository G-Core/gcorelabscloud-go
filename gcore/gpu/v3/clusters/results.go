package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster resource.
func (r commonResult) Extract() (*Cluster, error) {
	var s Cluster
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Image.
type GetResult struct {
	commonResult
}

type Interface struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Volume represents a volume structure.
type Volume struct {
	Size                 int                    `json:"size"`
	Source               string                 `json:"source"`
	Type                 string                 `json:"type"`
	DeletedOnTermination bool                   `json:"deleted_on_termination"`
	Metadata             map[string]interface{} `json:"metadata"`
	Name                 *string                `json:"name"`
	BootIndex            *int                   `json:"boot_index"`
	ImageID              *string                `json:"image_id"`
	SnapshotID           *string                `json:"snapshot_id"`
	VolumeID             *string                `json:"volume_id"`
}

type ClusterServerSettings struct {
	Interfaces       []Interface `json:"interfaces"`
	SecurityGroupIDs []string    `json:"security_groups"`
	Volumes          []Volume    `json:"volumes"`
	UserData         string      `json:"user_data"`
	KeypairID        *string     `json:"keypair_id"`
}

type Cluster struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Status          string                   `json:"status"`
	PlacementPolicy *string                  `json:"placement_policy"`
	FlavorID        string                   `json:"flavor_id"`
	ImageID         string                   `json:"image_id"`
	Metadata        map[string]interface{}   `json:"metadata"`
	ServersCount    int                      `json:"servers_count"`
	CreatedAt       gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt       *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	ServerIDs       *[]string                `json:"server_ids"`
	ServerSettings  ClusterServerSettings    `json:"server_settings"`
}
