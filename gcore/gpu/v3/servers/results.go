package servers

import (
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ServerPage is the page returned by a pager when traversing over a collection of servers.
type ServerPage struct {
	pagination.OffsetPageBase
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
