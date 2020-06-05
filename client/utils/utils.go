package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/fatih/structs"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var (
	slPfx = fmt.Sprintf("sl:::%d:::", time.Now().UTC().UnixNano())
)

type EnumValue struct {
	Enum        []string
	Default     string
	Destination *string
	selected    string
}

func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			if e.Destination != nil {
				e.Destination = &value
			}
			e.selected = value
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumValue) String() string {
	if e.selected == "" {
		return e.Default
	}
	return e.selected
}

type EnumStringSliceValue struct {
	Enum     []string
	Default  string
	selected []string
}

func (e *EnumStringSliceValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.selected = append(e.selected, value)
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumStringSliceValue) String() string {
	if len(e.selected) == 0 {
		return fmt.Sprintf("%s", []string{e.Default})
	}
	return fmt.Sprintf("%s", e.selected)
}

// Serialize allows EnumStringSliceValue to fulfill Serializer
func (e EnumStringSliceValue) Serialize() string {
	jsonBytes, _ := json.Marshal(e.selected)
	return fmt.Sprintf("%s%s", slPfx, string(jsonBytes))
}

// Value returns the slice of strings set by this flag
func (e EnumStringSliceValue) Value() []string {
	return e.selected
}

// Get returns the slice of strings set by this flag
func (e EnumStringSliceValue) Get() interface{} {
	return e
}

func GetEnumStringSliceValue(c *cli.Context, name string) []string {
	return c.Value(name).(EnumStringSliceValue).Value()
}

func tableHeaderFromStruct(m interface{}) []string {
	return structs.Names(m)
}

func tableRowFromStruct(m interface{}) []string {
	var res []string
	values := structs.Values(m)
	for _, v := range values {
		value, _ := json.Marshal(v)
		res = append(res, string(value))
	}
	return res
}

func renderTable(input interface{}) {
	if input == nil {
		return
	}
	results := interfaceToSlice(input)
	if len(results) == 0 {
		return
	}
	res := results[0]
	if reflect.TypeOf(res).Kind() != reflect.Struct {
		for i := 0; i < len(results); i++ {
			results[i] = struct {
				Value string
			}{Value: reflect.ValueOf(results[i]).String()}
		}
		res = results[0]
	}
	headers := tableHeaderFromStruct(res)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, res := range results {
		table.Append(tableRowFromStruct(res))
	}
	table.Render()
}

func interfaceToSlice(input interface{}) []interface{} {
	var records []interface{}
	if input == nil {
		return records
	}
	object := reflect.ValueOf(input)
	if reflect.TypeOf(input).Kind() != reflect.Slice {
		records = append(records, input)
		return records
	}
	for i := 0; i < object.Len(); i++ {
		records = append(records, object.Index(i).Interface())
	}
	return records
}

func renderJSON(input interface{}) error {
	if input == nil || (reflect.TypeOf(input).Kind() == reflect.Slice && reflect.ValueOf(input).Len() == 0) {
		return nil
	}
	res, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

func renderYAML(input interface{}) {
	if input == nil || (reflect.TypeOf(input).Kind() == reflect.Slice && reflect.ValueOf(input).Len() == 0) {
		return
	}
	res, _ := yaml.Marshal(input)
	fmt.Println(string(res))
}

func ShowResults(input interface{}, format string) {
	switch format {
	case "json":
		err := renderJSON(input)
		if err != nil {
			fmt.Println(err)
		}
	case "table":
		renderTable(input)
	case "yaml":
		renderYAML(input)
	}
}

func StringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func IntToPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func BoolToPointer(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func StringFromIndex(content []string, idx int, defaultValue string) string {
	if idx < len(content) {
		return content[idx]
	}
	return defaultValue
}

func IntFromIndex(content []int, idx int, defaultValue int) int {
	if idx < len(content) {
		return content[idx]
	}
	return defaultValue
}

func StringSliceToMap(slice []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong label format: %s", s)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

func StringSliceToMapNil(slice []string) (*map[string]string, error) {
	if len(slice) == 0 {
		return nil, nil
	}
	m, err := StringSliceToMap(slice)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func StringSliceToMapInterface(slice []string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong label format: %s", s)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

func WaitTaskAndShowResult(
	c *cli.Context,
	client *gcorecloud.ServiceClient,
	results *tasks.TaskResults,
	stopOnTaskError bool,
	infoRetriever tasks.RetrieveTaskResult,
) error {
	if c.Bool("wait") {
		if len(results.Tasks) == 0 {
			return cli.NewExitError(fmt.Errorf("wrong task response"), 1)
		}
		task := results.Tasks[0]
		waitSeconds := c.Int("wait-seconds")
		result, err := tasks.WaitTaskAndReturnResult(client, task, stopOnTaskError, waitSeconds, infoRetriever)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		ShowResults(result, c.String("format"))
	} else {
		ShowResults(results, c.String("format"))
	}
	return nil
}

func GetAbsPath(filename string) (string, error) {
	path, err := homedir.Expand(filename)
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func FileExists(filename string) (bool, error) {
	path, err := GetAbsPath(filename)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return !info.IsDir(), nil
}

func WriteToFile(filename string, content []byte) error {
	path, err := GetAbsPath(filename)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, content, 0644) // nolint
	return err
}

func ReadFile(filename string) ([]byte, error) {
	path, err := GetAbsPath(filename)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(path)
}

func CheckYamlFile(filename string) (content []byte, err error) {
	content, err = ReadFile(filename)
	if err != nil {
		return
	}
	out := make(map[string]interface{})
	err = yaml.Unmarshal(content, out)
	if err != nil {
		return
	}
	return
}

func getSliceLength(slice interface{}) (int, error) {
	switch t := slice.(type) {
	case []string:
		return len(t), nil
	case []int:
		return len(t), nil
	}
	return 0, fmt.Errorf("unexpected type: %T", slice)
}

func ValidateEqualSlicesLength(slices ...interface{}) error {
	if len(slices) == 0 || len(slices) == 1 {
		return nil
	}
	ln, err := getSliceLength(slices[0])
	if err != nil {
		return err
	}
	for _, slice := range slices[1:] {
		lns, err := getSliceLength(slice)
		if err != nil {
			return err
		}
		if lns != ln {
			return fmt.Errorf("element %v has different length than %v", slice, slices[0])
		}
	}
	return nil
}
