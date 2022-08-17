package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/baremetal/v1/bmcapacity"

const (
	availableNodesResponse = `
{
  "capacity": {
    "bm1-basic-small": 3,
    "bm1-infrastructure-small": 2,
    "bm1-infrastructure-medium": 1
  }
}
`
)

var (
	availableNodes = bmcapacity.AvailableNodes{
		Capacity: map[string]int{
			"bm1-basic-small":           3,
			"bm1-infrastructure-small":  2,
			"bm1-infrastructure-medium": 1,
		},
	}
)
