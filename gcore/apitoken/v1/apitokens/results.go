package apitokens

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/types"
)

type getResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts an api token resource.
func (r getResult) Extract() (*APIToken, error) {
	var s APIToken
	err := r.ExtractInto(&s)
	return &s, err
}

func (r getResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type listResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts an api token resources list.
func (r listResult) Extract() ([]APIToken, error) {
	var s []APIToken
	err := r.ExtractInto(&s)
	return s, err
}

func (r listResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "")
}

type createResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a token resource.
func (r createResult) Extract() (*Token, error) {
	var s Token
	err := r.ExtractInto(&s)
	return &s, err
}

func (r createResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as an APIToken.
type GetResult struct {
	getResult
}

// ListResult represents the result of a get operation. Call its Extract
// method to interpret it as an []APIToken.
type ListResult struct {
	listResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Token.
type CreateResult struct {
	createResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}

// APIToken represents an api token structure.
type APIToken struct {
	ID          int                         `json:"id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	ExpDate     *gcorecloud.JSONRFC3339ZZ   `json:"exp_date"`
	ClientUser  *ClientUser                 `json:"client_user"`
	Deleted     bool                        `json:"deleted"`
	Expired     bool                        `json:"expired"`
	Created     gcorecloud.JSONRFC3339ZZ    `json:"created"`
	LastUsage   *gcorecloud.JSONRFC3339Date `json:"last_usage"`
}

// Token represents a token structure.
type Token struct {
	Token string `json:"token"`
}

// ClientUser represents a client user structure.
type ClientUser struct {
	Role      ClientRole `json:"role"`
	Deleted   bool       `json:"deleted"`
	UserID    int        `json:"user_id"`
	UserName  string     `json:"user_name"`
	UserEmail string     `json:"user_email"`
	ClientID  int        `json:"client_id"`
}

// ClientRole represents a client role structure.
type ClientRole struct {
	ID   types.RoleIDType   `json:"id"`
	Name types.RoleNameType `json:"name"`
}
