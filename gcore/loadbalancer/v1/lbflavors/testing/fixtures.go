package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbflavors"

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "currency_code": "USD",
      "flavor_id": "g1-gpu-1-2-1",
      "flavor_name": "g1-gpu-1-2-1",
      "hardware_description": {
        "gpu": "1x NVIDIA 11GB"
      },
      "price_per_hour": 1,
      "price_per_month": 720,
      "price_status": "show",
      "ram": 2048,
      "vcpus": 1
    }
  ]
}
`

var (
	flavor1 = lbflavors.Flavor{
		PricePerMonth: 720,
		VCPUs:         1,
		FlavorName:    "g1-gpu-1-2-1",
		HardwareDescription: lbflavors.HardwareDescription{
			GPU: "1x NVIDIA 11GB",
		},
		CurrencyCode: "USD",
		PriceStatus:  "show",
		PricePerHour: 1,
		RAM:          2048,
		FlavorID:     "g1-gpu-1-2-1",
	}
	expectedFlavorList = []lbflavors.Flavor{flavor1}
)
