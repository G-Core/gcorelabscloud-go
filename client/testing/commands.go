package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"

	"github.com/G-Core/gcorelabscloud-go/cmd"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

func flagSet(name string, flags []cli.Flag) (*flag.FlagSet, error) {
	set := flag.NewFlagSet(name, flag.ContinueOnError)

	for _, f := range flags {
		if err := f.Apply(set); err != nil {
			return nil, err
		}
	}
	set.SetOutput(ioutil.Discard)
	return set, nil
}

func CheckMapInMap(srcMap map[string]interface{}, dstMap map[string]interface{}) bool {
	if len(srcMap) > len(dstMap) {
		return false
	}
	if len(srcMap) == len(dstMap) {
		return reflect.DeepEqual(srcMap, dstMap)
	}
	slicedMap := make(map[string]interface{}, len(srcMap))

	for k := range srcMap {
		if val, ok := dstMap[k]; ok {
			slicedMap[k] = val
		} else {
			return false
		}
	}

	return reflect.DeepEqual(srcMap, slicedMap)
}

func PrepareMetadata(apiMetadata []metadata.Metadata) (metadataMap map[string]string, metadataReadOnly []map[string]interface{}) {
	metadataMap = make(map[string]string)
	metadataReadOnly = make([]map[string]interface{}, 0, len(apiMetadata))

	if len(apiMetadata) > 0 {
		for _, metadataItem := range apiMetadata {
			if !metadataItem.ReadOnly {
				metadataMap[metadataItem.Key] = metadataItem.Value
			}
			metadataReadOnly = append(metadataReadOnly, map[string]interface{}{
				"key":       metadataItem.Key,
				"value":     metadataItem.Value,
				"read_only": metadataItem.ReadOnly,
			})
		}
	}

	return
}

func NormalizeMetadata(metadata interface{}, defaults ...bool) (map[string]interface{}, error) {
	normalizedMetadata := map[string]interface{}{}
	readOnly := false

	if len(defaults) > 0 {
		readOnly = defaults[0]
	}

	switch metadata.(type) {
	default:
		return nil, fmt.Errorf("unexpected type %T", metadata)
	case []map[string]interface{}:
		for _, v := range metadata.([]map[string]interface{}) {
			normalizedMetadata[v["key"].(string)] = v
		}
	case map[string]interface{}:
		for k, v := range metadata.(map[string]interface{}) {
			normalizedMetadata[k] = map[string]interface{}{
				"key":       k,
				"value":     v,
				"read_only": readOnly,
			}
		}
	case map[string]string:
		for k, v := range metadata.(map[string]string) {
			normalizedMetadata[k] = map[string]interface{}{
				"key":       k,
				"value":     v,
				"read_only": readOnly,
			}
		}
	}

	return normalizedMetadata, nil
}

func CompareMetadata(srcMeta interface{}, dstMeta interface{}) (bool, error) {
	src, err := NormalizeMetadata(srcMeta)
	if err != nil {
		return false, err
	}

	dst, err := NormalizeMetadata(dstMeta)
	if err != nil {
		return false, err
	}

	return CheckMapInMap(src, dst), nil
}
func InitTestApp(args []string) (*cli.App, *cli.Context) {
	a := cmd.NewApp(args)
	a.Setup()

	set, _ := flagSet(a.Name, a.Flags)
	ctx := cli.NewContext(a, set, &cli.Context{Context: context.Background()})
	return a, ctx
}

type FuncGetMetadata func() ([]metadata.Metadata, error)

func MetadataTest(getMetadata FuncGetMetadata, a *cli.App, resourceName string, resourceID string) error {
	// test metadata create
	var meta []map[string]interface{}
	args := []string{"gcoreclient", resourceName, "metadata", "create", "-m", "key1=val1", "-m", "key2=val2", resourceID}
	err := a.Run(args)
	if err != nil {
		return err
	}
	var typedMeta []metadata.Metadata

	typedMeta, err = getMetadata()
	if err != nil {
		return err
	}

	metadataMap, _ := PrepareMetadata(typedMeta)
	if !reflect.DeepEqual(metadataMap, map[string]string{"key1": "val1", "key2": "val2"}) {
		return errors.New("meta after create is not equal")
	}

	// test metadata list
	a.Writer = new(bytes.Buffer)
	args = []string{"gcoreclient", resourceName, "metadata", "list", "-t", resourceID}
	err = a.Run(args)
	if err != nil {
		return err
	}

	err = json.Unmarshal(a.Writer.(*bytes.Buffer).Bytes(), &meta)
	if err != nil {
		return err
	}

	isEqual, err := CompareMetadata(map[string]string{"key1": "val1", "key2": "val2"}, meta)

	if err != nil {
		return err
	}
	if !isEqual {
		return errors.New("meta after list is not equal")
	}

	// test metadata get by key
	var metaOne map[string]interface{}
	a.Writer = new(bytes.Buffer)
	args = []string{"gcoreclient", resourceName, "metadata", "show", "-t", "-m", "key1", resourceID}
	err = a.Run(args)
	if err != nil {
		return err
	}

	err = json.Unmarshal(a.Writer.(*bytes.Buffer).Bytes(), &metaOne)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(map[string]interface{}{"key": "key1", "value": "val1", "read_only": false}, metaOne) {
		return errors.New("meta after show is not equal")
	}

	// test metadata update
	args = []string{"gcoreclient", resourceName, "metadata", "update", "-m", "key1=val11", "-m", "key4=val4", resourceID}
	err = a.Run(args)
	if err != nil {
		return err
	}

	typedMeta, err = getMetadata()
	if err != nil {
		return err
	}
	metadataMap, _ = PrepareMetadata(typedMeta)
	if !reflect.DeepEqual(metadataMap, map[string]string{"key1": "val11", "key2": "val2", "key4": "val4"}) {
		return errors.New("metadata after update not equal")
	}

	// test metadata replace
	args = []string{"gcoreclient", resourceName, "metadata", "replace", "-m", "key1=val11", "-m", "key4=val4", resourceID}
	err = a.Run(args)
	if err != nil {
		return err
	}

	typedMeta, err = getMetadata()
	if err != nil {
		return err
	}
	metadataMap, _ = PrepareMetadata(typedMeta)
	if !reflect.DeepEqual(metadataMap, map[string]string{"key1": "val11", "key4": "val4"}) {
		return errors.New("metadata after replace not equal")
	}

	return nil
}

const networkDeleting int = 1200
const networkCreatingTimeout int = 1200

func CreateTestNetwork(client *gcorecloud.ServiceClient, opts networks.CreateOpts) (string, error) {
	res, err := networks.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	networkID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, networkCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		networkID, err := networks.ExtractNetworkIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve network ID from task info: %w", err)
		}
		return networkID, nil
	},
	)
	if err != nil {
		return "", err
	}
	return networkID.(string), nil
}

func DeleteTestNetwork(client *gcorecloud.ServiceClient, networkID string) error {
	results, err := networks.Delete(client, networkID).Extract()
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, networkDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := networks.Get(client, networkID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete network with ID: %s", networkID)
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
