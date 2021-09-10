package schedules

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
)

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToScheduleUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a schedule.
type UpdateOpts struct {
	Hour                 string                         `json:"hour,omitempty"`
	Month                string                         `json:"month,omitempty"`
	Timezone             string                         `json:"timezone,omitempty"`
	Weeks                int                            `json:"weeks,omitempty"`
	Minute               string                         `json:"minute,omitempty"`
	Day                  string                         `json:"day,omitempty"`
	RetentionTime        lifecyclepolicy.RetentionTimer `json:"retention_time"`
	Time                 string                         `json:"time,omitempty"`
	ResourceNameTemplate string                         `json:"resource_name_template,omitempty"`
	Days                 int                            `json:"days,omitempty"`
	Minutes              int                            `json:"minutes,omitempty"`
	Week                 string                         `json:"week,omitempty"`
	MaxQuantity          int                            `json:"max_quantity" required:"true"`
	DayOfWeek            string                         `json:"day_of_week,omitempty"`
	Hours                int                            `json:"hours,omitempty"`
	Type                 lifecyclepolicy.ScheduleType   `json:"type,omitempty"`
}

func (opts UpdateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToScheduleUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToScheduleUpdateMap() (map[string]interface{}, error) {
	err := gcorecloud.TranslateValidationError(opts.Validate())
	if err != nil {
		return nil, err
	}

	body, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, err
}

// Get retrieves a specific schedule based on its unique ID.
func Get(c *gcorecloud.ServiceClient, scheduleID string) (r GetResult) {
	url := resourceURL(c, scheduleID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Update accepts a UpdateOpts struct and updates an existing schedule using the
// values provided.
func Update(c *gcorecloud.ServiceClient, scheduleID string, opts UpdateOptsBuilder) (r GetResult) {
	b, err := opts.ToScheduleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(resourceURL(c, scheduleID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
