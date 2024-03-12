package volumes

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

type VolumeTaskResult struct {
	Volumes []string `json:"volumes"`
}

func ExtractVolumeIDFromTask(task *tasks.Task) (string, error) {
	var result VolumeTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode volume information in task structure: %w", err)
	}
	if len(result.Volumes) == 0 {
		return "", fmt.Errorf("cannot decode volume information in task structure: %w", err)
	}
	return result.Volumes[0], nil
}
