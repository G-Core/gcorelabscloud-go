package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

const allowedAddressPairsRequest = `
{
  "allowed_address_pairs": [
    {
      "ip_address": "192.168.123.20",
      "mac_address": "00:16:3e:f2:87:16"
    }
  ]
}
`

const allowedAddressPairsResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	PortID     = "1f0ca628-a73b-42c0-bdac-7b10d023e097"
	PortIPRaw1 = "192.168.123.20"
	Tasks1     = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
)
