package testing

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/client/images/v1/client"
	gtest "github.com/G-Core/gcorelabscloud-go/client/testing"
	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

func TestImageMetadata(t *testing.T) {
	resourceName := "image"
	args := []string{"gcoreclient", resourceName}
	a, ctx := gtest.InitTestApp(args)

	cli, err := client.NewImageClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	downloadClient, err := client.NewDownloadImageClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := images.UploadOpts{
		HwMachineType:  "q35",
		SshKey:         "allow",
		Name:           "test_image_tf1",
		OSType:         "linux",
		URL:            "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img",
		HwFirmwareType: "uefi",
		Metadata:       map[string]string{"key1": "val1", "key2": "val2"},
		Architecture:   "x86_64",
	}

	res, err := images.Upload(downloadClient, opts).Extract()

	if err != nil {
		t.Fatal(err)
	}

	taskID := res.Tasks[0]
	resourceID, err := tasks.WaitTaskAndReturnResult(downloadClient, taskID, true, 1200, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(downloadClient, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		resourceID, err := images.ExtractImageIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve volume ID from task info: %w", err)
		}
		return resourceID, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	typedResourceID := resourceID.(string)
	defer images.Delete(cli, typedResourceID)

	err = gtest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := images.Get(cli, typedResourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, typedResourceID)

	if err != nil {
		t.Fatal(err)
	}
}
