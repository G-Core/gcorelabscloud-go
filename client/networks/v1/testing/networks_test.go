package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/client/networks/v1/client"
	gtest "github.com/G-Core/gcorelabscloud-go/client/testing"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

func TestNetworksMetadata(t *testing.T) {
	resourceName := "network"
	args := []string{"gcoreclient", resourceName}
	a, ctx := gtest.InitTestApp(args)

	resourceClient, err := client.NewNetworkClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := networks.CreateOpts{
		Name: "test-network1",
	}

	resourceID, err := gtest.CreateTestNetwork(resourceClient, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer gtest.DeleteTestNetwork(resourceClient, resourceID)

	err = gtest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := networks.Get(resourceClient, resourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
