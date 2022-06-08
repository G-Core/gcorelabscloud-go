package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const getNamespaceResponse = `
{
  "name": "namespace-name",
  "description": "Namespace description",
  "envs": {
    "ENV_VAR": "value 1"
  }
}
`

const listNamespaceResponse = `
{
  "count": 1,
  "results": [
    {
      "name": "namespace-name",
      "description": "Namespace description",
      "envs": {
        "ENV_VAR": "value 1"
      }
    }
  ]
}
`

const createNamespaceRequest = `
{
  "name": "string",
  "description": "long string",
  "envs": {
    "ENV_VAR": "value 1"
  }
}
`

const updateNamespaceRequest = `
{
  "description": "long string",
  "envs": {
    "ENV_VAR": "value 1"
  }
}
`

const listFunctionResponse = `
{
  "count": 1,
  "results": [
    {
      "name": "function-name",
      "description": "Function description",
      "envs": {
        "ENV_VAR": "value 1",
        "ENVIRONMENT_VARIABLE": "value 2"
      },
      "runtime": "python3.7.12",
      "timeout": 5,
      "flavor": "64mCPU-64MB",
      "autoscaling": {
        "min_instances": 1,
        "max_instances": 2
      },
      "code_text": "def main(): print('It works!')",
      "main_method": "main"
    }
  ]
}
`

const getFunctionResponse = `
{
  "name": "function-name",
  "description": "Function description",
  "envs": {
	"ENV_VAR": "value 1",
	"ENVIRONMENT_VARIABLE": "value 2"
  },
  "runtime": "python3.7.12",
  "timeout": 5,
  "flavor": "64mCPU-64MB",
  "autoscaling": {
	"min_instances": 1,
	"max_instances": 2
  },
  "code_text": "def main(): print('It works!')",
  "main_method": "main"
}
`
const createFunctionRequest = `
{
  "name": "function-name",
  "description": "Function description",
  "envs": {
    "ENV_VAR": "value 1"
  },
  "runtime": "python3.7.12",
  "timeout": 5,
  "flavor": "64mCPU-64MB",
  "autoscaling": {
    "min_instances": 1,
    "max_instances": 2
  },
  "code_text": "def main(): print('It works!')",
  "main_method": "main"
}
`
const updateFunctionRequest = `
{
  "description": "string",
  "envs": {
    "property1": "string"
  },
  "code_text": "string",
  "timeout": 180,
  "autoscaling": {
    "min_instances": 1,
    "max_instances": 2
  },
  "main_method": "string",
  "flavor": "string"
}
`

const taskResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	nsName     = "test-namespace"
	expectedNs = faas.Namespace{
		Name:        "namespace-name",
		Description: "Namespace description",
		Envs: map[string]string{
			"ENV_VAR": "value 1",
		},
	}
	expectedNsSlice = []faas.Namespace{expectedNs}
	tasks1          = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	fName     = "function-name"
	expectedF = faas.Function{
		Name:        fName,
		Description: "Function description",
		Envs: map[string]string{
			"ENV_VAR":              "value 1",
			"ENVIRONMENT_VARIABLE": "value 2",
		},
		Runtime: "python3.7.12",
		Timeout: 5,
		Flavor:  "64mCPU-64MB",
		Autoscaling: faas.FunctionAutoscaling{
			MinInstances: 1,
			MaxInstances: 2,
		},
		CodeText:   "def main(): print('It works!')",
		MainMethod: "main",
	}
	expectedFSlice = []faas.Function{expectedF}
)
