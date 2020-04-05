package snapshots

import (
	"fmt"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a snapshot resource.
func (r commonResult) Extract() (*Snapshot, error) {
	var s Snapshot
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTasks is a function that accepts a result and extracts a snapshot creation task resource.
func (r commonResult) ExtractTasks() (*tasks.TaskResults, error) {
	var t tasks.TaskResults
	err := r.ExtractInto(&t)
	return &t, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Snapshot.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Snapshot.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// Snapshot represents a snapshot.
type Snapshot struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	Status        string                   `json:"status"`
	Size          int                      `json:"size"`
	VolumeID      string                   `json:"volume_id"`
	CreatedAt     gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt     *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	Metadata      map[string]interface{}   `json:"metadata"`
	CreatorTaskID *string                  `json:"creator_task_id"`
	TaskID        *string                  `json:"task_id"`
	ProjectID     int                      `json:"project_id"`
	RegionID      int                      `json:"region_id"`
	Region        string                   `json:"region"`
}

// SnapshotPage is the page returned by a pager when traversing over a
// collection of snapshots.
type SnapshotPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of snapshots has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SnapshotPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SnapshotPage struct is empty.
func (r SnapshotPage) IsEmpty() (bool, error) {
	is, err := ExtractSnapshots(r)
	return len(is) == 0, err
}

// ExtractSnapshot accepts a Page struct, specifically a SnapshotPage struct,
// and extracts the elements into a slice of Snapshot structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSnapshots(r pagination.Page) ([]Snapshot, error) {
	var s []Snapshot
	err := ExtractSnapshotsInto(r, &s)
	return s, err
}

func ExtractSnapshotsInto(r pagination.Page, v interface{}) error {
	return r.(SnapshotPage).Result.ExtractIntoSlicePtr(v, "results")
}

type SnapshotTaskResult struct {
	Snapshots []string `json:"snapshots"`
}

func ExtractSnapshotIDFromTask(task *tasks.Task) (string, error) {
	var result SnapshotTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode snapshot information in task structure: %w", err)
	}
	if len(result.Snapshots) == 0 {
		return "", fmt.Errorf("cannot decode snapshot information in task structure: %w", err)
	}
	return result.Snapshots[0], nil
}
