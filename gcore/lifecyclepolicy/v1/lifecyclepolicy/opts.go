package lifecyclepolicy

import "github.com/G-Core/gcorelabscloud-go"

// Options used for query parameters

// GetOpts represents options for Get.
type GetOpts struct {
	NeedVolumes bool `q:"need_volumes"`
}

// ListOpts represents options for List.
type ListOpts GetOpts

// Options used for request body

func ValidateAndBuildRequestBody(opts interface{}) (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateOpts represents options for Create.
type CreateOpts struct {
	Name      string               `json:"name" validate:"required,name"`
	Status    PolicyStatus         `json:"status,omitempty" validate:"omitempty,enum"`
	Action    PolicyAction         `json:"action" validate:"required,enum"`
	Schedules []CreateScheduleOpts `json:"schedules,omitempty" validate:"dive"`
	VolumeIds []string             `json:"volume_ids,omitempty"`
}

// UpdateOpts represents options for Update.
type UpdateOpts struct {
	Name   string       `json:"name,omitempty" validate:"omitempty,name"`
	Status PolicyStatus `json:"status,omitempty" validate:"omitempty,enum"`
}

// AddVolumesOpts represents options for AddVolumes.
// Volumes already managed by policy are ignored.
type AddVolumesOpts struct {
	VolumeIds []string `json:"volume_ids" validate:"required"`
}

// RemoveVolumesOpts represents options for RemoveVolumes.
type RemoveVolumesOpts struct {
	VolumeIds []string `json:"volume_ids" validate:"required"`
}

// CreateScheduleOpts represents options used to create a single schedule.
type CreateScheduleOpts interface {
	SetCommonCreateScheduleOpts(opts CommonCreateScheduleOpts)
}

type CommonCreateScheduleOpts struct {
	Type                 ScheduleType    `json:"type" validate:"required,enum"`
	ResourceNameTemplate string          `json:"resource_name_template,omitempty"`
	MaxQuantity          int             `json:"max_quantity" validate:"required,gt=0,lt=10001"`
	RetentionTime        *RetentionTimer `json:"retention_time,omitempty"`
}

// CreateCronScheduleOpts represents options used to create a single cron schedule.
type CreateCronScheduleOpts struct { // TODO: validate?
	CommonCreateScheduleOpts
	Timezone  string `json:"timezone,omitempty"`
	Week      string `json:"week,omitempty"`
	DayOfWeek string `json:"day_of_week,omitempty"`
	Month     string `json:"month,omitempty"`
	Day       string `json:"day,omitempty"`
	Hour      string `json:"hour,omitempty"`
	Minute    string `json:"minute,omitempty"`
}

// CreateIntervalScheduleOpts represents options used to create a single interval schedule.
type CreateIntervalScheduleOpts struct {
	CommonCreateScheduleOpts
	Weeks   int `json:"weeks,omitempty" validate:"required_without_all=Days Hours Minutes"`
	Days    int `json:"days,omitempty" validate:"required_without_all=Weeks Hours Minutes"`
	Hours   int `json:"hours,omitempty" validate:"required_without_all=Weeks Days Minutes"`
	Minutes int `json:"minutes,omitempty" validate:"required_without_all=Weeks Days Hours"`
}

func (opts *CreateCronScheduleOpts) SetCommonCreateScheduleOpts(common CommonCreateScheduleOpts) {
	opts.CommonCreateScheduleOpts = common
}
func (opts *CreateIntervalScheduleOpts) SetCommonCreateScheduleOpts(common CommonCreateScheduleOpts) {
	opts.CommonCreateScheduleOpts = common
}

// AddSchedulesOpts represents options for AddSchedules.
type AddSchedulesOpts struct {
	Schedules []CreateScheduleOpts `json:"schedules" validate:"required,dive"`
}

// RemoveSchedulesOpts represents options for RemoveSchedules.
type RemoveSchedulesOpts struct {
	ScheduleIDs []string `json:"schedule_ids" validate:"required"`
}
