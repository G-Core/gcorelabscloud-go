package testing

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/loadbalancers/v1/client"
	gtest "github.com/G-Core/gcorelabscloud-go/client/testing"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

const (
	LoadBalancersPoint        = "loadbalancers"
	LoadBalancerCreateTimeout = 2400
	lbTestName                = "test-lb"
	lbListenerTestName        = "test-listener"
)

func createTestLoadBalancerWithListener(client *gcorecloud.ServiceClient, opts loadbalancers.CreateOpts) (string, error) {
	res, err := loadbalancers.Create(client, opts, nil).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	lbID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		lbID, err := loadbalancers.ExtractLoadBalancerIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return lbID, nil
	})
	if err != nil {
		return "", err
	}
	return lbID.(string), nil
}

func TestLBSMetadata(t *testing.T) {
	resourceName := "loadbalancer"
	args := []string{"gcoreclient", resourceName}
	a, ctx := gtest.InitTestApp(args)

	client, err := client.NewLoadbalancerClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := loadbalancers.CreateOpts{
		Name: lbTestName,
		Listeners: []loadbalancers.CreateListenerOpts{{
			Name:         lbListenerTestName,
			ProtocolPort: 80,
			Protocol:     types.ProtocolTypeHTTP,
		}},
	}

	resourceID, err := createTestLoadBalancerWithListener(client, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer loadbalancers.Delete(client, resourceID, nil)

	err = gtest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := loadbalancers.Get(client, resourceID, nil).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
