package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	gtest "github.com/G-Core/gcorelabscloud-go/client/testing"
	"github.com/G-Core/gcorelabscloud-go/client/volumes/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"reflect"
	"testing"
)

func TestVolumeMetadata(t *testing.T) {
	args := []string{"gcoreclient", "volume"}
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
	volumeID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, 1200, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		volumeID, err := volumes.ExtractVolumeIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve volume ID from task info: %w", err)
		}
		return volumeID, nil
	})

	if err != nil {
		t.Fatal(err)
	}
	volID := volumeID.(string)

	// test metadata create
	var meta []map[string]interface{}
	args = []string{"gcoreclient", "volume", "metadata", "create", "-m", "key1=val1", "-m", "key2=val2", volID}
	a.Run(args)

	volume, err := volumes.Get(client, volID).Extract()
	if err != nil {
		t.Fatal(err)
	}

	metadataMap, _ := gtest.PrepareMetadata(volume.Metadata)
	if !reflect.DeepEqual(metadataMap, map[string]string{"key1": "val1", "key2": "val2"}) {
		t.Fatal("metadata not equal")
	}

	// test metadata list
	a.Writer = new(bytes.Buffer)
	args = []string{"gcoreclient", "volume", "metadata", "list", "-t", volID}
	a.Run(args)

	err = json.Unmarshal(a.Writer.(*bytes.Buffer).Bytes(), &meta)
	if err != nil {
		t.Fatal(err)
	}

	isEqual, err := gtest.CompareMetadata(map[string]string{"key1": "val1", "key2": "val2"}, meta)

	if err != nil {
		t.Fatal(err)
	}

	if !isEqual {
		t.Fatal("meta is not equal")
	}

	// test metadata get by key
	var metaOne map[string]interface{}
	a.Writer = new(bytes.Buffer)
	args = []string{"gcoreclient", "volume", "metadata", "show", "-t", "-m", "key1", volID}
	a.Run(args)

	err = json.Unmarshal(a.Writer.(*bytes.Buffer).Bytes(), &metaOne)
	if err != nil {
		t.Fatal(err)
	}
	isEqual, err = gtest.CompareMetadata(map[string]string{"key1": "val1"}, metaOne)

	if err != nil {
		t.Fatal(err)
	}

	// test metadata update
	args = []string{"gcoreclient", "volume", "metadata", "update", "-m", "key1=val11", "-m", "key4=val4", volID}
	a.Run(args)

	volume, err = volumes.Get(client, volID).Extract()
	if err != nil {
		t.Fatal(err)
	}
	metadataMap, _ = gtest.PrepareMetadata(volume.Metadata)
	if !reflect.DeepEqual(metadataMap, map[string]string{"key1": "val11", "key2": "val2", "key4": "val4"}) {
		t.Fatal("metadata not equal")
	}

	// test metadata replace
	args = []string{"gcoreclient", "volume", "metadata", "replace", "-m", "key1=val11", "-m", "key4=val4", volID}
	a.Run(args)

	volume, err = volumes.Get(client, volID).Extract()
	if err != nil {
		t.Fatal(err)
	}
	metadataMap, _ = gtest.PrepareMetadata(volume.Metadata)
	if !reflect.DeepEqual(metadataMap, map[string]string{"key1": "val11", "key4": "val4"}) {
		t.Fatal("metadata not equal")
	}

	defer volumes.Delete(client, volID, volumes.DeleteOpts{})
}
