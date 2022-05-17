package testing

import "fmt"

var (
	createRequests = []string{
		`
{
  "name": "TestPolicy0",
  "action": "volume_snapshot",
  "status": "active",
  "volume_ids": [
    "c2d7afb7-888c-4234-8da0-6c3fc9298c17",
	"c646322c-c4b8-4f04-83cc-fa9c68a280c2"
  ],
  "schedules": [
	{
      "type": "interval",
      "max_quantity": 4,
      "retention_time": {},
	  "resource_name_template": "test template name {volume_id}",
	  "hours": 60,
	  "minutes": 80,
	  "weeks": 20,
	  "days": 40
	},
	{
	  "type": "interval",
	  "max_quantity": 5,
	  "resource_name_template": "another {volume_id} test template name",
	  "days": 1
	},
    {
      "type": "cron",
      "max_quantity": 5,
      "resource_name_template": "yet another test template name",
      "retention_time": {
        "days": 40
      },
      "timezone": "Europe/Kirov",
      "week": "1,3,5",
      "day_of_week": "mon,sun,3",
      "month": "2,3,4",
      "day": "*",
      "hour": "17,19",
      "minute": "13"
    }
  ]
}`,
		`
{
  "name": "TestPolicy1",
  "action": "volume_snapshot",
  "status": "paused",
  "volume_ids": [
	"c2d7afb7-888c-4234-8da0-6c3fc9298c17"
  ]
}`,
	}
	updateRequests = []string{
		`
{
  "name": "TestPolicy0",
  "status": "active"
}`,
		`
{
  "name": "TestPolicy1",
  "status": "paused"
}`,
	}
	addVolumesRequests = []string{
		`{"volume_ids": ["c2d7afb7-888c-4234-8da0-6c3fc9298c17", "c646322c-c4b8-4f04-83cc-fa9c68a280c2"]}`,
		`{"volume_ids": ["c2d7afb7-888c-4234-8da0-6c3fc9298c17"]}`,
	}
	removeVolumesRequests = addVolumesRequests
	addSchedulesRequests  = `
{
  "schedules": [
	{
      "type": "interval",
      "max_quantity": 4,
      "retention_time": {},
	  "resource_name_template": "test template name {volume_id}",
	  "hours": 60,
	  "minutes": 80,
	  "weeks": 20,
	  "days": 40
	},
	{
	  "type": "interval",
	  "max_quantity": 5,
	  "resource_name_template": "another {volume_id} test template name",
	  "days": 1
	},
    {
      "type": "cron",
      "max_quantity": 5,
      "resource_name_template": "yet another test template name",
      "retention_time": {
        "days": 40
      },
      "timezone": "Europe/Kirov",
      "week": "1,3,5",
      "day_of_week": "mon,sun,3",
      "month": "2,3,4",
      "day": "*",
      "hour": "17,19",
      "minute": "13"
    }
  ]
}`
	removeSchedulesRequest = `
{
  "schedule_ids": [
	"aa14fdc6-6589-4136-be94-275456df9d79",
	"3a03c757-1f24-43a0-a5b4-942f326f7d30",
	"47aee364-ad88-4b66-aed2-6ed8f635ee17"
  ]
}`
	responsesWithVolumes = []string{
		`
{
  "name": "TestPolicy0",
  "id": 113,
  "action": "volume_snapshot",
  "project_id": 1,
  "status": "active",
  "user_id": 4585,
  "region_id": 1,
  "volumes": [
    {
      "volume_id": "c2d7afb7-888c-4234-8da0-6c3fc9298c17",
      "volume_name": "Vol_0"
    },
    {
      "volume_id": "c646322c-c4b8-4f04-83cc-fa9c68a280c2",
      "volume_name": "Vol_1"
    }
  ],
  "schedules": [
	{
      "user_id": 4585,
      "type": "interval",
      "owner": "lifecycle_policy",
      "max_quantity": 4,
      "retention_time": {
		"minutes": 0,
		"days": 0,
		"hours": 0,
		"weeks": 0
	  },
	  "owner_id": 113,
	  "resource_name_template": "test template name {volume_id}",
	  "id": "aa14fdc6-6589-4136-be94-275456df9d79",
	  "hours": 60,
	  "minutes": 80,
	  "weeks": 20,
	  "days": 40
	},
	{
	  "user_id": 4585,
	  "type": "interval",
	  "owner": "lifecycle_policy",
	  "max_quantity": 5,
	  "owner_id": 113,
	  "resource_name_template": "another {volume_id} test template name",
	  "id": "3a03c757-1f24-43a0-a5b4-942f326f7d30",
	  "hours": 0,
	  "minutes": 0,
	  "weeks": 0,
	  "days": 1
	},
    {
      "type": "cron",
      "id": "47aee364-ad88-4b66-aed2-6ed8f635ee17",
      "owner": "lifecycle_policy",
      "owner_id": 113,
      "max_quantity": 5,
      "user_id": 4585,
      "resource_name_template": "yet another test template name",
      "retention_time": {
		"minutes": 0,
        "days": 40,
		"hours": 0,
		"weeks": 0
      },
      "timezone": "Europe/Kirov",
      "week": "1,3,5",
      "day_of_week": "mon,sun,3",
      "month": "2,3,4",
      "day": "*",
      "hour": "17,19",
      "minute": "13"
    }
  ]
}`,
		`
{
  "name": "TestPolicy1",
  "id": 128,
  "action": "volume_snapshot",
  "project_id": 1,
  "status": "paused",
  "user_id": 4585,
  "region_id": 1,
  "volumes": [
    {
      "volume_id": "c2d7afb7-888c-4234-8da0-6c3fc9298c17",
      "volume_name": "Vol_0"
    }
  ],
  "schedules": []
}`,
	}
	responsesWithoutVolumes = []string{
		`
{
  "name": "TestPolicy0",
  "id": 113,
  "action": "volume_snapshot",
  "project_id": 1,
  "status": "active",
  "user_id": 4585,
  "region_id": 1,
  "volumes": null,
  "schedules": [
	{
      "user_id": 4585,
      "type": "interval",
      "owner": "lifecycle_policy",
      "max_quantity": 4,
      "retention_time": {
		"minutes": 0,
		"days": 0,
		"hours": 0,
		"weeks": 0
	  },
	  "owner_id": 113,
	  "resource_name_template": "test template name {volume_id}",
	  "id": "aa14fdc6-6589-4136-be94-275456df9d79",
	  "hours": 60,
	  "minutes": 80,
	  "weeks": 20,
	  "days": 40
	},
	{
	  "user_id": 4585,
	  "type": "interval",
	  "owner": "lifecycle_policy",
	  "max_quantity": 5,
	  "owner_id": 113,
	  "resource_name_template": "another {volume_id} test template name",
	  "id": "3a03c757-1f24-43a0-a5b4-942f326f7d30",
	  "hours": 0,
	  "minutes": 0,
	  "weeks": 0,
	  "days": 1
	},
    {
      "type": "cron",
      "id": "47aee364-ad88-4b66-aed2-6ed8f635ee17",
      "owner": "lifecycle_policy",
      "owner_id": 113,
      "max_quantity": 5,
      "user_id": 4585,
      "resource_name_template": "yet another test template name",
      "retention_time": {
		"minutes": 0,
        "days": 40,
		"hours": 0,
		"weeks": 0
      },
      "timezone": "Europe/Kirov",
      "week": "1,3,5",
      "day_of_week": "mon,sun,3",
      "month": "2,3,4",
      "day": "*",
      "hour": "17,19",
      "minute": "13"
    }
  ]
}`,
		`
{
  "name": "TestPolicy1",
  "id": 128,
  "action": "volume_snapshot",
  "project_id": 1,
  "status": "paused",
  "user_id": 4585,
  "region_id": 1,
  "volumes": null,
  "schedules": []
}`,
	}
	listPoliciesResponseWithVolumes = fmt.Sprintf(`{"count":2,"results":[%s,%s]}`, responsesWithVolumes[0], responsesWithVolumes[1])
)
