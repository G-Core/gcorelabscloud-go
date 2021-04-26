package lifecyclepolicy

import (
	"encoding/json"
	"fmt"
)

type (
	ScheduleType string
	PolicyStatus string
	PolicyAction string
)

const (
	PolicyTypeCron             ScheduleType = "cron"
	PolicyTypeInterval         ScheduleType = "interval"
	PolicyStatusActive         PolicyStatus = "active"
	PolicyStatusPaused         PolicyStatus = "paused"
	PolicyActionVolumeSnapshot PolicyAction = "volume_snapshot"
)

func (t ScheduleType) List() []ScheduleType {
	return []ScheduleType{PolicyTypeInterval, PolicyTypeCron}
}

func (t ScheduleType) String() string {
	return string(t)
}

func (t ScheduleType) StringList() []string {
	var strings []string
	for _, x := range t.List() {
		strings = append(strings, x.String())
	}
	return strings
}

func (t ScheduleType) IsValid() error {
	for _, x := range t.List() {
		if t == x {
			return nil
		}
	}
	return fmt.Errorf("invalid schedule type: %v", t)
}

func (s PolicyStatus) List() []PolicyStatus {
	return []PolicyStatus{PolicyStatusPaused, PolicyStatusActive}
}

func (s PolicyStatus) String() string {
	return string(s)
}

func (s PolicyStatus) StringList() []string {
	var strings []string
	for _, x := range s.List() {
		strings = append(strings, x.String())
	}
	return strings
}

func (s PolicyStatus) IsValid() error {
	for _, x := range s.List() {
		if s == x {
			return nil
		}
	}
	return fmt.Errorf("invalid lifecycle policy status: %v", s)
}

func (a PolicyAction) List() []PolicyAction {
	return []PolicyAction{PolicyActionVolumeSnapshot}
}

func (a PolicyAction) String() string {
	return string(a)
}

func (a PolicyAction) StringList() []string {
	var strings []string
	for _, x := range a.List() {
		strings = append(strings, x.String())
	}
	return strings
}

func (a PolicyAction) IsValid() error {
	for _, x := range a.List() {
		if a == x {
			return nil
		}
	}
	return fmt.Errorf("invalid lifecycle policy action: %v", a)
}

// Schedule represents a schedule resource.
type Schedule interface {
	GetCommonSchedule() CommonSchedule
}

type RetentionTimer struct {
	Weeks   int `json:"weeks"`
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type CommonSchedule struct {
	Type                 ScheduleType    `json:"type"`
	ID                   string          `json:"id"`
	Owner                string          `json:"owner"`
	OwnerID              int             `json:"owner_id"`
	MaxQuantity          int             `json:"max_quantity"`
	UserID               int             `json:"user_id"`
	ResourceNameTemplate string          `json:"resource_name_template"`
	RetentionTime        *RetentionTimer `json:"retention_time"`
}

type IntervalSchedule struct {
	CommonSchedule
	Weeks   int `json:"weeks"`
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type CronSchedule struct {
	CommonSchedule
	Timezone  string `json:"timezone"`
	Week      string `json:"week"`
	DayOfWeek string `json:"day_of_week"`
	Month     string `json:"month"`
	Day       string `json:"day"`
	Hour      string `json:"hour"`
	Minute    string `json:"minute"`
}

func (s CronSchedule) GetCommonSchedule() CommonSchedule {
	return s.CommonSchedule
}
func (s IntervalSchedule) GetCommonSchedule() CommonSchedule {
	return s.CommonSchedule
}

// rawSchedule is internal struct for unmarshalling into Schedule.
type rawSchedule struct {
	json.RawMessage
}

// cook is internal method for unmarshalling rawSchedule into Schedule
func (r rawSchedule) cook() (Schedule, error) {
	var typeStruct struct {
		ScheduleType `json:"type"`
	}
	if err := json.Unmarshal(r.RawMessage, &typeStruct); err != nil {
		return nil, err
	}
	switch typeStruct.ScheduleType {
	default:
		return nil, fmt.Errorf("unexpected schedule type %s", typeStruct.ScheduleType)
	case PolicyTypeCron:
		var cronSchedule CronSchedule
		if err := json.Unmarshal(r.RawMessage, &cronSchedule); err != nil {
			return nil, err
		}
		return cronSchedule, nil
	case PolicyTypeInterval:
		var intervalSchedule IntervalSchedule
		if err := json.Unmarshal(r.RawMessage, &intervalSchedule); err != nil {
			return nil, err
		}
		return intervalSchedule, nil
	}
}

// Volume represents a volume resource.
type Volume struct {
	ID   string `json:"volume_id"`
	Name string `json:"volume_name"`
}

// LifecyclePolicy represents a lifecycle policy resource.
type LifecyclePolicy struct {
	Name      string       `json:"name"`
	ID        int          `json:"id"`
	Action    PolicyAction `json:"action"`
	ProjectID int          `json:"project_id"`
	Status    PolicyStatus `json:"status"`
	UserID    int          `json:"user_id"`
	RegionID  int          `json:"region_id"`
	Volumes   []Volume     `json:"volumes"`
	Schedules []Schedule   `json:"schedules"`
}

// rawLifecyclePolicy is internal struct for unmarshalling into LifecyclePolicy.
type rawLifecyclePolicy struct {
	Name      string        `json:"name"`
	ID        int           `json:"id"`
	Action    PolicyAction  `json:"action"`
	ProjectID int           `json:"project_id"`
	Status    PolicyStatus  `json:"status"`
	UserID    int           `json:"user_id"`
	RegionID  int           `json:"region_id"`
	Volumes   []Volume      `json:"volumes"`
	Schedules []rawSchedule `json:"schedules"`
}

// cook is internal method for unmarshalling rawLifecyclePolicy into LifecyclePolicy
func (rawPolicy rawLifecyclePolicy) cook() (*LifecyclePolicy, error) {
	schedules := make([]Schedule, len(rawPolicy.Schedules))
	for i, b := range rawPolicy.Schedules {
		s, err := b.cook()
		if err != nil {
			return nil, err
		}
		schedules[i] = s
	}
	rawPolicy.Schedules = nil
	b, err := json.Marshal(rawPolicy)
	if err != nil {
		return nil, err
	}
	var policy LifecyclePolicy
	if err := json.Unmarshal(b, &policy); err != nil {
		return nil, err
	}
	policy.Schedules = schedules
	return &policy, nil
}
