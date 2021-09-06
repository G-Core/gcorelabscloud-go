package quotas

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

type Quota map[string]int

type CombinedQuota struct {
	GlobalQuotas   Quota   `json:"global_quotas"`
	RegionalQuotas []Quota `json:"regional_quotas"`
}

type CommonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a quotas resource.
func (r CommonResult) Extract() (*Quota, error) {
	var s Quota
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CommonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoMapPtr(v, "")
}

type CombinedResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a combined quota resource.
func (r CombinedResult) Extract() (*CombinedQuota, error) {
	var s CombinedQuota
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CombinedResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}
