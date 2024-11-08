package testing

import (
	"fmt"
	metadataV1 "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"
	metadataV2 "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v2/metadata"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	"testing"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	v1client "github.com/G-Core/gcorelabscloud-go/client/loadbalancers/v1/client"
	v2client "github.com/G-Core/gcorelabscloud-go/client/loadbalancers/v2/client"
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

func createMetadataV2(clientV1 *gcorecloud.ServiceClient, clientV2 *gcorecloud.ServiceClient, lbId string, opts map[string]string) (interface{}, error) {
	res, err := metadataV2.MetadataCreateOrUpdate(clientV2, lbId, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	state, err := tasks.WaitTaskAndReturnResult(clientV1, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientV1, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		state := taskInfo.State
		if state != tasks.TaskStateFinished {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return state, nil
	})
	if err != nil {
		return "", err
	}
	return state, nil
}

func replaceMetadataV2(clientV1 *gcorecloud.ServiceClient, clientV2 *gcorecloud.ServiceClient, lbId string, opts map[string]string) (interface{}, error) {
	res, err := metadataV2.MetadataReplace(clientV2, lbId, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	state, err := tasks.WaitTaskAndReturnResult(clientV1, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientV1, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		state := taskInfo.State
		if state != tasks.TaskStateFinished {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return state, nil
	})
	if err != nil {
		return "", err
	}
	return state, nil
}

func deleteMetadataV2(clientV1 *gcorecloud.ServiceClient, clientV2 *gcorecloud.ServiceClient, lbId string, key string) (interface{}, error) {
	res, err := metadataV2.MetadataDelete(clientV2, lbId, key).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	state, err := tasks.WaitTaskAndReturnResult(clientV1, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientV1, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		state := taskInfo.State
		if state != tasks.TaskStateFinished {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return state, nil
	})
	if err != nil {
		return "", err
	}
	return state, nil
}

func TestLBSMetadataV2(t *testing.T) {
	resourceName := "loadbalancer"
	args := []string{"gcoreclient", resourceName}
	_, ctx := gtest.InitTestApp(args)

	clientV1, err := v1client.NewLoadbalancerClientV1(ctx)
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

	resourceID, err := createTestLoadBalancerWithListener(clientV1, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer loadbalancers.Delete(clientV1, resourceID, nil)

	clientV2, err := v2client.NewLoadbalancerClientV2(ctx)
	if err != nil {
		t.Fatal(err)
	}

	createMetadataOpts := map[string]string{
		"test1": "test1",
		"test2": "test2",
	}
	createMetadataTaskState, err := createMetadataV2(clientV1, clientV2, resourceID, createMetadataOpts)
	if createMetadataTaskState != tasks.TaskStateFinished {
		t.Fatal(err)
	}

	createMetadataListAll, err := metadataV1.MetadataListAll(clientV1, resourceID)
	if err != nil {
		t.Fatal(err)
	}

	test1Metadata := metadata.Metadata{
		Key:      "test1",
		Value:    "test1",
		ReadOnly: false,
	}
	test2Metadata := metadata.Metadata{
		Key:      "test2",
		Value:    "test2",
		ReadOnly: false,
	}
	createExpectedMetadataList := []metadata.Metadata{test1Metadata, test2Metadata}
	th.AssertDeepEquals(t, createExpectedMetadataList, createMetadataListAll)

	time.Sleep(3 * time.Second) // need few more seconds while LB will be Active

	deleteMetadataTaskState, err := deleteMetadataV2(clientV1, clientV2, resourceID, test1Metadata.Key)
	if deleteMetadataTaskState != tasks.TaskStateFinished {
		t.Fatal(err)
	}

	deleteMetadataListAll, err := metadataV1.MetadataListAll(clientV1, resourceID)
	if err != nil {
		t.Fatal(err)
	}

	deleteExpectedMetadataList := []metadata.Metadata{test2Metadata}
	th.AssertDeepEquals(t, deleteExpectedMetadataList, deleteMetadataListAll)

	replaceMetadataOpts := map[string]string{
		"123": "321",
	}

	time.Sleep(3 * time.Second) // need few more seconds while LB will be Active

	replaceMetadataTaskState, err := replaceMetadataV2(clientV1, clientV2, resourceID, replaceMetadataOpts)
	if replaceMetadataTaskState != tasks.TaskStateFinished {
		t.Fatal(err)
	}

	replaceMetadataListAll, err := metadataV1.MetadataListAll(clientV1, resourceID)
	if err != nil {
		t.Fatal(err)
	}

	onetwothreeMetadata := metadata.Metadata{
		Key:      "123",
		Value:    "321",
		ReadOnly: false,
	}
	updateExpectedMetadataList := []metadata.Metadata{onetwothreeMetadata}
	th.AssertDeepEquals(t, updateExpectedMetadataList, replaceMetadataListAll)
}
