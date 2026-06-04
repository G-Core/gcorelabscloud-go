package pagination

import (
	"strconv"
)

const (
	defaultOffset = 0
)

type listResult struct {
	Count   int           `json:"count"`
	Results []interface{} `json:"results"`
}

// OffsetPageBase may be embedded to implement a page that operates on offset / limit query parameters.
type OffsetPageBase struct {
	PageResult
}

// NextPageURL constructs next page URL using offset / limit query parameters.
func (current OffsetPageBase) NextPageURL() (string, error) {
	var res listResult
	if err := current.Result.ExtractInto(&res); err != nil {
		return "", err
	}

	// Advance by the number of results actually returned on this page rather
	// than by an assumed page size. When the request carried no explicit limit
	// the API returns the whole collection in a single page, so there is no
	// next page to fetch.
	got := len(res.Results)
	if got == 0 {
		return "", nil
	}

	offset, err := current.getQueryParam("offset", defaultOffset)
	if err != nil {
		return "", err
	}

	nextOffset := offset + got
	if nextOffset >= res.Count {
		return "", nil
	}

	query := current.URL.Query()
	query.Set("offset", strconv.Itoa(nextOffset))

	nextURL := current.URL
	nextURL.RawQuery = query.Encode()
	return nextURL.String(), nil
}

// IsEmpty returns true when the page is empty, otherwise false.
func (current OffsetPageBase) IsEmpty() (bool, error) {
	var res listResult
	err := current.Result.ExtractInto(&res)
	return len(res.Results) == 0, err
}

// GetBody returns the page's body.
func (current OffsetPageBase) GetBody() interface{} {
	return current.Body
}

func (current OffsetPageBase) getQueryParam(name string, def int) (int, error) {
	if current.Query().Get(name) == "" {
		return def, nil
	}
	v, err := strconv.Atoi(current.Query().Get(name))
	if err != nil {
		return def, err
	}
	return v, nil
}
