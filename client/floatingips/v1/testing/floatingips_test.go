package testing

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/floatingips/v1/client"
	gtest "github.com/G-Core/gcorelabscloud-go/client/testing"
	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

const (
	FloatingIPCreateTimeout = 1200
)

func createTestFloatingIP(client *gcorecloud.ServiceClient, opts floatingips.CreateOpts) (string, error) {
	res, err := floatingips.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	floatingIPID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, FloatingIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		floatingIPID, err := floatingips.ExtractFloatingIPIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve FloatingIP ID from task info: %w", err)
		}
		return floatingIPID, nil
	})

	if err != nil {
		return "", err
	}
	return floatingIPID.(string), nil
}

func TestFipsMetadata(t *testing.T) {
	resourceName := "floatingip"
	args := []string{"gcoreclient", resourceName}
	a, ctx := gtest.InitTestApp(args)

	client, err := client.NewFloatingIPClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := floatingips.CreateOpts{}

	resourceID, err := createTestFloatingIP(client, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer floatingips.Delete(client, resourceID)

	err = gtest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := floatingips.Get(client, resourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
