package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/flavors"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "name": "inference-16vcpu-232gib-1xh100-80gb",
      "cpu": 2,
      "memory": 4,
      "gpu": 1,
      "gpu_model": "H100",
      "gpu_memory": 80,
      "is_gpu_shared": false,
      "gpu_compute_capability": "8.6"
    }
  ]
}
`

const GetResponse = `
{
  "name": "inference-16vcpu-232gib-1xh100-80gb",
  "cpu": 2,
  "memory": 4,
  "gpu": 1,
  "gpu_model": "H100",
  "gpu_memory": 80,
  "is_gpu_shared": false,
  "gpu_compute_capability": "8.6"
}
`

var (
	Flavor1 = flavors.Flavor{
		Name:                 "inference-16vcpu-232gib-1xh100-80gb",
		Cpu:                  2,
		Memory:               4,
		Gpu:                  1,
		GpuModel:             "H100",
		GpuMemory:            80,
		IsGpuShared:          false,
		GpuComputeCapability: "8.6",
	}
	FlavorSlice = []flavors.Flavor{Flavor1}
)
