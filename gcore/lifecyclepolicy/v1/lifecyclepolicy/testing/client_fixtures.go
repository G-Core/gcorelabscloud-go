package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

var (
	volumes = []lifecyclepolicy.Volume{
		{
			ID:   "c2d7afb7-888c-4234-8da0-6c3fc9298c17",
			Name: "Vol_0",
		},
		{
			ID:   "c646322c-c4b8-4f04-83cc-fa9c68a280c2",
			Name: "Vol_1",
		},
	}
	schedules = []lifecyclepolicy.Schedule{
		lifecyclepolicy.IntervalSchedule{
			Weeks:   20,
			Days:    40,
			Hours:   60,
			Minutes: 80,
			CommonSchedule: lifecyclepolicy.CommonSchedule{
				Type:                 lifecyclepolicy.PolicyTypeInterval,
				ID:                   "aa14fdc6-6589-4136-be94-275456df9d79",
				Owner:                "lifecycle_policy",
				OwnerID:              113,
				MaxQuantity:          4,
				UserID:               4585,
				ResourceNameTemplate: "test template name {volume_id}",
				RetentionTime:        &lifecyclepolicy.RetentionTimer{},
			},
		},
		lifecyclepolicy.IntervalSchedule{
			Weeks:   0,
			Days:    1,
			Hours:   0,
			Minutes: 0,
			CommonSchedule: lifecyclepolicy.CommonSchedule{
				Type:                 lifecyclepolicy.PolicyTypeInterval,
				ID:                   "3a03c757-1f24-43a0-a5b4-942f326f7d30",
				Owner:                "lifecycle_policy",
				OwnerID:              113,
				MaxQuantity:          5,
				UserID:               4585,
				ResourceNameTemplate: "another {volume_id} test template name",
				RetentionTime:        nil,
			},
		},
		lifecyclepolicy.CronSchedule{
			Timezone:  "Europe/Kirov",
			Week:      "1,3,5",
			DayOfWeek: "mon,sun,3",
			Month:     "2,3,4",
			Day:       "*",
			Hour:      "17,19",
			Minute:    "13",
			CommonSchedule: lifecyclepolicy.CommonSchedule{
				Type:                 lifecyclepolicy.PolicyTypeCron,
				ID:                   "47aee364-ad88-4b66-aed2-6ed8f635ee17",
				Owner:                "lifecycle_policy",
				OwnerID:              113,
				MaxQuantity:          5,
				UserID:               4585,
				ResourceNameTemplate: "yet another test template name",
				RetentionTime: &lifecyclepolicy.RetentionTimer{
					Days: 40,
				},
			},
		},
	}
	policies = []lifecyclepolicy.LifecyclePolicy{
		{
			Name:      "TestPolicy0",
			ID:        113,
			Action:    lifecyclepolicy.PolicyActionVolumeSnapshot,
			ProjectID: fake.ProjectID,
			Status:    lifecyclepolicy.PolicyStatusActive,
			UserID:    4585,
			RegionID:  fake.RegionID,
			Volumes:   []lifecyclepolicy.Volume{volumes[0], volumes[1]},
			Schedules: schedules,
		},
		{
			Name:      "TestPolicy1",
			ID:        128,
			Action:    lifecyclepolicy.PolicyActionVolumeSnapshot,
			ProjectID: fake.ProjectID,
			Status:    lifecyclepolicy.PolicyStatusPaused,
			UserID:    4585,
			RegionID:  fake.RegionID,
			Volumes:   []lifecyclepolicy.Volume{volumes[0]},
			Schedules: []lifecyclepolicy.Schedule{},
		},
	}
	createScheduleOpts = []lifecyclepolicy.CreateScheduleOpts{
		&lifecyclepolicy.CreateIntervalScheduleOpts{
			Weeks:   20,
			Days:    40,
			Hours:   60,
			Minutes: 80,
			CommonCreateScheduleOpts: lifecyclepolicy.CommonCreateScheduleOpts{
				Type:                 lifecyclepolicy.PolicyTypeInterval,
				MaxQuantity:          4,
				ResourceNameTemplate: "test template name {volume_id}",
				RetentionTime:        &lifecyclepolicy.RetentionTimer{},
			},
		},
		&lifecyclepolicy.CreateIntervalScheduleOpts{
			Days: 1,
			CommonCreateScheduleOpts: lifecyclepolicy.CommonCreateScheduleOpts{
				Type:                 "interval",
				ResourceNameTemplate: "another {volume_id} test template name",
				MaxQuantity:          5,
			},
		},
		&lifecyclepolicy.CreateCronScheduleOpts{
			Timezone:  "Europe/Kirov",
			Week:      "1,3,5",
			DayOfWeek: "mon,sun,3",
			Month:     "2,3,4",
			Day:       "*",
			Hour:      "17,19",
			Minute:    "13",
			CommonCreateScheduleOpts: lifecyclepolicy.CommonCreateScheduleOpts{
				Type:                 lifecyclepolicy.PolicyTypeCron,
				MaxQuantity:          5,
				ResourceNameTemplate: "yet another test template name",
				RetentionTime: &lifecyclepolicy.RetentionTimer{
					Days: 40,
				},
			},
		},
	}
	createOpts = []lifecyclepolicy.CreateOpts{
		{
			Name:      "TestPolicy0",
			Action:    lifecyclepolicy.PolicyActionVolumeSnapshot,
			Status:    lifecyclepolicy.PolicyStatusActive,
			VolumeIds: []string{volumes[0].ID, volumes[1].ID},
			Schedules: createScheduleOpts,
		},
		{
			Name:      "TestPolicy1",
			Action:    lifecyclepolicy.PolicyActionVolumeSnapshot,
			Status:    lifecyclepolicy.PolicyStatusPaused,
			VolumeIds: []string{volumes[0].ID},
			Schedules: nil,
		},
	}
	updateOpts = []lifecyclepolicy.UpdateOpts{
		{
			Name:   "TestPolicy0",
			Status: lifecyclepolicy.PolicyStatusActive,
		},
		{
			Name:   "TestPolicy1",
			Status: lifecyclepolicy.PolicyStatusPaused,
		},
	}
	addVolumesOpts = []lifecyclepolicy.AddVolumesOpts{
		{
			VolumeIds: []string{volumes[0].ID, volumes[1].ID},
		},
		{
			VolumeIds: []string{volumes[0].ID},
		},
	}
	removeVolumesOpts = []lifecyclepolicy.RemoveVolumesOpts{
		{
			VolumeIds: []string{volumes[0].ID, volumes[1].ID},
		},
		{
			VolumeIds: []string{volumes[0].ID},
		},
	}
	addSchedulesOpts = lifecyclepolicy.AddSchedulesOpts{
		Schedules: createScheduleOpts,
	}
	removeSchedulesOpts = lifecyclepolicy.RemoveSchedulesOpts{
		ScheduleIDs: []string{
			schedules[0].GetCommonSchedule().ID,
			schedules[1].GetCommonSchedule().ID,
			schedules[2].GetCommonSchedule().ID,
		},
	}
)
