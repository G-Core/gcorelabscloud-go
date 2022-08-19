package testing

import (
	"context"
	"flag"
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/cmd"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"reflect"
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
