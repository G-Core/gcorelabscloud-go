package tasks

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	"time"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToTaskListQuery() (string, error)
}

// ListOpts allows the filtering and sorting List API response.
type ListOpts struct {
	ProjectID     *int       `q:"project_id"`
	State         *TaskState `q:"state"`
	TaskType      *string    `q:"task_type"`
	Sorting       *string    `q:"sorting"`
	FromTimestamp *string    `q:"from_timestamp"`
}

// isValid validation for Task options.
func (opts ListOpts) isValid() error {
	if opts.State != nil {
		switch TaskState(*opts.State) {
		case TaskStateError, TaskStateFinished, TaskStateNew, TaskStateRunning:
			fallthrough
		default:
			return fmt.Errorf(`invalid task type: "%s"`, opts.TaskType)
		}
	}

	if opts.Sorting != nil {
		switch *opts.Sorting {
		case TaskSortingOldFirst, TaskSortingToNewFirst:
			fallthrough
		default:
			return fmt.Errorf(`invalid task sort option: "%s"`, opts.Sorting)
		}
	}

	if opts.FromTimestamp != nil {
		_, err := time.Parse(time.RFC3339, *opts.FromTimestamp)
		if err != nil {
			return fmt.Errorf(`timestamp "%s" should be ISO format string`, opts.FromTimestamp)
		}
	}

	return nil
}

// ToTaskListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTaskListQuery() (string, error) {
	if err := opts.isValid(); err != nil {
		return "", fmt.Errorf(`invalid task list filter params: "%s"`, opts)
	}

	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// cluster templates. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToTaskListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return TaskPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll returns all Tasks.
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Task, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractTasks(page)
}

// Get retrieves a specific cluster template based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
