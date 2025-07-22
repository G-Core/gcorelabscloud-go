package pagination

import (
	"strconv"
)

const (
	defaultOffset = 0
	defaultLimit  = 10
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
	offset, err := current.getQueryParam("offset", defaultOffset)
	if err != nil {
		return "", err
	}

	limit, err := current.getQueryParam("limit", defaultLimit)
	if err != nil {
		return "", err
	}

	query := current.URL.Query()
	query.Set("offset", strconv.Itoa(offset+limit))
	query.Set("limit", strconv.Itoa(limit))

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
