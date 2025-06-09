package servers

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ServerPage is the page returned by a pager when traversing over a collection of servers.
type ServerPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of AI Clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ServerPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ServerPage struct is empty.
func (r ServerPage) IsEmpty() (bool, error) {
	is, err := ExtractServers(r)
	return len(is) == 0, err
}

// ExtractServers accepts a Page struct, specifically a ServerPage struct,
// and extracts the elements into a slice of Server structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractServers(r pagination.Page) ([]Server, error) {
	var s []Server
	err := ExtractServersInto(r, &s)
	return s, err
}

func ExtractServersInto(r pagination.Page, v interface{}) error {
	return r.(ServerPage).Result.ExtractIntoSlicePtr(v, "results")
}
