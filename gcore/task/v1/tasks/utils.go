package tasks

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// WaitForStatus will continually poll the task resource, checking for a particular
// status. It will do this for the amount of seconds defined.
func WaitForStatus(client *gcorecloud.ServiceClient, id string, status TaskState, secs int, stopOnTaskError bool) error {
	return gcorecloud.WaitFor(secs, func() (bool, error) {
		task, err := Get(client, id).Extract()
		if err != nil {
			return false, err
		}

		if task.State == status {
			return true, nil
		}

		if task.State == TaskStateError {
			errorText := ""
			if task.Error != nil {
				errorText = *task.Error
			}
			return false, fmt.Errorf("task is in error state: %s. Error: %s", task.State, errorText)
		}

		if task.Error != nil && stopOnTaskError {
			return false, fmt.Errorf("task is in error state: %s", *task.Error)
		}

		return false, nil
	})
}

// WaitTaskAndProcessResult periodically check status state and invoke taskProcessor when when task is finished
func WaitTaskAndProcessResult(
	client *gcorecloud.ServiceClient, task TaskID, stopOnTaskError bool, waitSeconds int, taskProcessor CheckTaskResult) error {
	err := WaitForStatus(client, string(task), TaskStateFinished, waitSeconds, stopOnTaskError)
	if err != nil {
		return err
	}
	err = taskProcessor(task)
	if err != nil {
		return err
	}
	return nil
}

// WaitTaskAndReturnResult periodically check status state and return changed object when task is finished
func WaitTaskAndReturnResult(
	client *gcorecloud.ServiceClient, task TaskID, stopOnTaskError bool,
	waitSeconds int, taskProcessor RetrieveTaskResult) (interface{}, error) {

	err := WaitForStatus(client, string(task), TaskStateFinished, waitSeconds, stopOnTaskError)
	if err != nil {
		return nil, err
	}
	result, err := taskProcessor(task)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type RetrieveTaskResult func(task TaskID) (interface{}, error)
type CheckTaskResult func(task TaskID) error
