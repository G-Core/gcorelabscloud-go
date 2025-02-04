package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/aiflavors"
)

const ListResponse = `
{
    "count": 1,
    "results": [
        {
            "resource_class": "bm1-ai-small",
            "hardware_description": {
                "network": "2x100G",
                "ipu": "vPOD-16 (Classic)",
                "poplar_count": 2
            },
            "disabled": false,
            "flavor_name": "bm1-ai-2xsmall-v1pod-16",
            "flavor_id": "bm1-ai-2xsmall-v1pod-16"
        }
	]
}
`

var (
	AIFlavor1 = aiflavors.AIFlavor{
		FlavorID:      "bm1-ai-2xsmall-v1pod-16",
		FlavorName:    "bm1-ai-2xsmall-v1pod-16",
		Disabled:      false,
		ResourceClass: "bm1-ai-small",
		HardwareDescription: &aiflavors.HardwareDescription{
			Network:     "2x100G",
			IPU:         "vPOD-16 (Classic)",
			PoplarCount: 2,
		},
	}
	ExpectedAIFlavorSlice = []aiflavors.AIFlavor{AIFlavor1}
)
