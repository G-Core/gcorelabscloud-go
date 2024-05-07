package volumes

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

type volumeAttachTaskResult struct {
	Volume_ID string `json:"volume_id"`
}

func ExtractVolumeIDFromAttachTask(task *tasks.Task) (string, error) {
	var result volumeAttachTaskResult
	err := gcorecloud.NativeMapToStruct(task.Data, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode volume information in task structure: %w", err)
	}
	if result.Volume_ID == "" {
		return "", fmt.Errorf("cannot decode volume information in task structure: %w", err)
	}
	return result.Volume_ID, nil
}
