package volumes

import (
	"time"

	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	pagination.LinkedPageBase
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a volume resource.
func (r commonResult) Extract() (*ClusterInstanceVolumesList, error) {
	var s ClusterInstanceVolumesList
	err := r.ExtractInto(&s)
	return &s, err
}

// ClusterInstanceVolumesList represents a paginated list of volumes
type ClusterInstanceVolumesList struct {
	Count   int                     `json:"count"`
	Results []ClusterInstanceVolume `json:"results"`
}

// GetBody implements pagination.Page interface
func (r ClusterInstanceVolumesList) GetBody() interface{} {
	return r
}

// IsEmpty checks if the volume list is empty
func (r commonResult) IsEmpty() (bool, error) {
	volumes, err := ExtractVolumes(r)
	return len(volumes) == 0, err
}

// ExtractVolumes extracts volumes from a ClusterInstanceVolumesList
func ExtractVolumes(r pagination.Page) ([]ClusterInstanceVolume, error) {
	var s ClusterInstanceVolumesList
	err := r.(commonResult).ExtractInto(&s)
	return s.Results, err
}

// DetailedMetadata represents metadata with read-only flag
type DetailedMetadata struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

// ClusterInstanceVolume represents a volume in a cluster
type ClusterInstanceVolume struct {
	ID        string             `json:"id"`
	ServerID  string             `json:"server_id"`
	Name      string             `json:"name"`
	Type      string             `json:"type"`
	Status    string             `json:"status"`
	Size      int                `json:"size"`
	Bootable  bool               `json:"bootable"`
	Metadata  []DetailedMetadata `json:"metadata"`
	CreatedAt time.Time          `json:"created_at"`
}
