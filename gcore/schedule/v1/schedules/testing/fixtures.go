package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"

const GetResponseCron = `
{
  "owner": "lifecycle_policy",
  "owner_id": 1,
  "id": "1488e2ce-f906-47fb-ba32-c25a3f63df4f",
  "max_quantity": 2,
  "day_of_week": "fri",
  "hour": "0, 20",
  "minute": "30",
  "type": "cron",
  "resource_name_template": "reserve snap of the volume {volume_id}",
  "timezone": "Europe/London",
  "retention_time": {
    "weeks": 2
  },
  "user_id": 12
}
`

const GetResponseInterval = `
{
  "owner": "lifecycle_policy",
  "owner_id": 1,
  "id": "1488e2ce-f906-47fb-ba32-c25a3f63df4f",
  "max_quantity": 2,
  "type": "interval",
  "resource_name_template": "reserve snap of the volume {volume_id}",
  "retention_time": {
    "weeks": 2
  },
  "user_id": 12,
  "weeks": 1,
  "days": 2,
  "hours": 3,
  "minutes": 4
}
`

const UpdateRequest = `
{
  "max_quantity": 2,
  "day_of_week": "fri, tue",
  "hour": "0, 20",
  "minute": "30",
  "type": "cron",
  "resource_name_template": "reserve snap of the volume {volume_id}",
  "timezone": "Europe/London",
  "retention_time": {
    "weeks": 2
  }
}
`

var (
	ScheduleCron1 = lifecyclepolicy.CronSchedule{
		CommonSchedule: lifecyclepolicy.CommonSchedule{
			Type:                 lifecyclepolicy.ScheduleTypeCron,
			ID:                   "1488e2ce-f906-47fb-ba32-c25a3f63df4f",
			Owner:                "lifecycle_policy",
			OwnerID:              1,
			MaxQuantity:          2,
			UserID:               12,
			ResourceNameTemplate: "reserve snap of the volume {volume_id}",
			RetentionTime: &lifecyclepolicy.RetentionTimer{
				Weeks: 2,
			},
		},
		Timezone:  "Europe/London",
		DayOfWeek: "fri",
		Hour:      "0, 20",
		Minute:    "30",
	}
	ScheduleInterval1 = lifecyclepolicy.IntervalSchedule{
		CommonSchedule: lifecyclepolicy.CommonSchedule{
			Type:                 lifecyclepolicy.ScheduleTypeInterval,
			ID:                   "1488e2ce-f906-47fb-ba32-c25a3f63df4f",
			Owner:                "lifecycle_policy",
			OwnerID:              1,
			MaxQuantity:          2,
			UserID:               12,
			ResourceNameTemplate: "reserve snap of the volume {volume_id}",
			RetentionTime: &lifecyclepolicy.RetentionTimer{
				Weeks: 2,
			},
		},
		Weeks:   1,
		Days:    2,
		Hours:   3,
		Minutes: 4,
	}
)
