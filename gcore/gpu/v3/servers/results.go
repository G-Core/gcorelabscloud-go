package servers

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// ListResult represents the result of a List operation.
type ListResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a list of servers.
func (r ListResult) Extract() (*ServersPage, error) {
	var s ServersPage
	err := r.ExtractInto(&s)
	return &s, err
}

// ServersPage represents a page of servers results.
type ServersPage struct {
	Count   int      `json:"count"`
	Results []Server `json:"results"`
}

// IsEmpty returns true if a ServersPage contains no results.
func (r ServersPage) IsEmpty() bool {
	return len(r.Results) == 0
}

// ExtractServers accepts a ListResult and extracts the servers.
func ExtractServers(lr ListResult) ([]Server, error) {
	var s ServersPage
	err := lr.ExtractInto(&s)
	return s.Results, err
}
