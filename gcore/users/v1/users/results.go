package users

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

type commonResult struct {
	gcorecloud.Result
}

type User struct {
	UserID int `json:"user_id"`
}

// Extract is a function that accepts a result and extracts a user resource.
func (r commonResult) Extract() (*User, error) {
	var s User
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a user.
type CreateResult struct {
	commonResult
}
