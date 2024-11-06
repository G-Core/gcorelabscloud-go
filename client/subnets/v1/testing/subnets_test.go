package testing

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/subnets/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"net"
	"testing"

	netclient "github.com/G-Core/gcorelabscloud-go/client/networks/v1/client"
	gtest "github.com/G-Core/gcorelabscloud-go/client/testing"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

const SubnetDeleting int = 1200
const SubnetCreatingTimeout int = 1200

func createTestSubnet(client *gcorecloud.ServiceClient, opts subnets.CreateOpts, subCidr string) (string, error) {
	var gccidr gcorecloud.CIDR
	_, netIPNet, err := net.ParseCIDR(subCidr)
	if err != nil {
		return "", err
	}
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask
	opts.CIDR = gccidr

	res, err := subnets.Create(client, opts, nil).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	subnetID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, SubnetCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Subnet, err := subnets.ExtractSubnetIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Subnet ID from task info: %w", err)
		}
		return Subnet, nil
	},
	)

	return subnetID.(string), err
}

func deleteTestSubnet(client *gcorecloud.ServiceClient, subnetID string) error {
	results, err := subnets.Delete(client, subnetID, nil).Extract()
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, SubnetDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := subnets.Get(client, subnetID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete subnet with ID: %s", subnetID)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	return err
}

func TestSubnetsMetadata(t *testing.T) {
	resourceName := "subnet"

	args := []string{"gcoreclient", resourceName}
	a, ctx := gtest.InitTestApp(args)

	netClient, err := netclient.NewNetworkClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	resourceClient, err := client.NewSubnetClientV1(ctx)

	opts := networks.CreateOpts{
		Name: "test-network1",
	}

	networkID, err := gtest.CreateTestNetwork(netClient, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer gtest.DeleteTestNetwork(netClient, networkID)

	optsSubnet := subnets.CreateOpts{
		Name:      "test-subnet",
		NetworkID: networkID,
	}

	resourceID, err := createTestSubnet(resourceClient, optsSubnet, "192.168.42.0/24")
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestSubnet(resourceClient, resourceID)

	err = gtest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := subnets.Get(resourceClient, resourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
