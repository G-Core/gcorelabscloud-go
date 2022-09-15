package testing

import (
	"fmt"
	"testing"

	gtest "github.com/G-Core/gcorelabscloud-go/client/testing"
	"github.com/G-Core/gcorelabscloud-go/client/volumes/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
)

func TestVolumeMetadata(t *testing.T) {
	resourceName := "volume"
	args := []string{"gcoreclient", resourceName}
	a, ctx := gtest.InitTestApp(args)

	client, err := client.NewVolumeClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := volumes.CreateOpts{
		Name:     "test-volume-1",
		Size:     1,
		Source:   volumes.NewVolume,
		TypeName: volumes.Standard,
	}

	res, err := volumes.Create(client, opts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	taskID := res.Tasks[0]
	resourceID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, 1200, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		resourceID, err := volumes.ExtractVolumeIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve volume ID from task info: %w", err)
		}
		return resourceID, nil
	})

	if err != nil {
		t.Fatal(err)
	}
	typedResourceID := resourceID.(string)
	defer volumes.Delete(client, typedResourceID, volumes.DeleteOpts{})

	err = gtest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := volumes.Get(client, typedResourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, typedResourceID)

	if err != nil {
		t.Fatal(err)
	}

}
