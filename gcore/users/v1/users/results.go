package users

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

type commonUserResult struct {
	gcorecloud.Result
}

type User struct {
	UserID int `json:"user_id"`
}

// Extract is a function that accepts a result and extracts a user resource.
func (r commonUserResult) Extract() (*User, error) {
	var s User
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonUserResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateUserResult represents the result of a create operation. Call its Extract
// method to interpret it as a user.
type CreateUserResult struct {
	commonUserResult
}

type commonApiTokenResult struct {
	gcorecloud.Result
}

type ApiToken struct {
	Token string `json:"token"`
}

// Extract is a function that accepts a result and extracts an api token resource.
func (r commonApiTokenResult) Extract() (*ApiToken, error) {
	var s ApiToken
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonApiTokenResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateApiTokenResult represents the result of a create operation. Call its Extract
// method to interpret it as an api token.
type CreateApiTokenResult struct {
	commonApiTokenResult
}
