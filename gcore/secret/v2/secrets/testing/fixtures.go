package testing

import (
	"fmt"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

var CreateRequest = fmt.Sprintf(`
{
  "name": "%s",
  "expiration": "%s",
  "payload": {
    "certificate": "%s",
    "private_key": "%s",
    "certificate_chain": "%s"
  }
}
`, secretName, createdTime, cert, privateKey, certChain)

var CreateResponse = `
{
  "tasks": [
    "d478ae29-dedc-4869-82f0-96104425f565"
  ]
}
`

var (
	createdTime          = "2025-12-28T19:14:44.180394"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339MilliNoZ, createdTime)
	cert                 = `a`
	privateKey           = `b`
	certChain            = `a`

	secretName = "AES key"
	Tasks1     = tasks.TaskResults{
		Tasks: []tasks.TaskID{"d478ae29-dedc-4869-82f0-96104425f565"},
	}
)
