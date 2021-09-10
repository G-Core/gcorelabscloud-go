package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/regionaccess/v1/regionsaccess"

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "access_all_edge_regions": true,
      "client_id": 8,
      "id": 11,
      "region_ids": [
        1,
        2,
        3
      ],
      "reseller_id": 1
    }
  ]
}
`

const CreateResponse = `
{
  "access_all_edge_regions": true,
  "client_id": 8,
  "region_ids": [
    1,
    2,
    3
  ],
  "reseller_id": 1
}
`

const CreateRequest = `
{
  "access_all_edge_regions": true,
  "client_id": 8,
  "region_ids": [
    1,
    2,
    3
  ],
  "reseller_id": 1
}
`

var (
	clientID      = 8
	resellerID    = 1
	RegionAccess1 = regionsaccess.RegionAccess{
		ID:                   11,
		AccessAllEdgeRegions: true,
		ClientID:             &clientID,
		RegionIDs:            []int{1, 2, 3},
		ResellerID:           &resellerID,
	}
	ExpectedRegionAccessSlice = []regionsaccess.RegionAccess{RegionAccess1}

	RegionAccessCreated1 = regionsaccess.RegionAccess{
		AccessAllEdgeRegions: true,
		ClientID:             &clientID,
		RegionIDs:            []int{1, 2, 3},
		ResellerID:           &resellerID,
	}
)
