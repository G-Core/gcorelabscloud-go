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
	  "enable_api_key": true,
	  "keys"    :["key-one"],
	  "disabled": false,
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
  "enable_api_key": true,
  "keys"    :["key-one"],
  "disabled": false,
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
  "enable_api_key": true,
  "keys"   :["key-one"],
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

const getKeyResponse = `
{
  "description": "description",
  "name": "test-key",
  "status": "active",
  "functions": [
    {
        "name": "function",
        "namespace": "namespace"
    }
  ]
}
`

const createKeyRequest = `
{
  "description": "description",
  "name": "test-key",
  "functions": [
    {
        "name": "function",
        "namespace": "namespace"
    }
  ]
}
`

const createKeyResponse = `
{
  "description": "description",
  "name": "test-key",
  "status": "active",
  "functions": [
    {
        "name": "function",
        "namespace": "namespace"
    }
  ]
}
`

const listKeysResponse = `
{
  "count": 1,
  "results": [
	  {
		  "description": "description",
		  "name": "test-key",
		  "status": "active",
		  "functions": [
			{
				"name": "function",
				"namespace": "namespace"
			}
		  ]
	  }
  ]
}
`

const updateKeyRequest = `
{
  "description": "long string",
  "functions": [
    {
      "name":      "function1",
      "namespace": "namespace1"
    },
    {
      "name":      "function2",
      "namespace": "namespace1"
    }
  ]
}
`

const updateKeyResponse = `
{
	"description":  "long string",
	"name": "test-key",
	"status": "active",
	"functions": [
		{
			"name":      "function1",
			"namespace": "namespace1"
		},
		{
			"name":      "function2",
			"namespace": "namespace1"
		}
	]
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
	min       = 1
	max       = 2
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
			MinInstances: &min,
			MaxInstances: &max,
		},
		CodeText:     "def main(): print('It works!')",
		MainMethod:   "main",
		EnableAPIKey: true,
		Disabled:     false,
		Keys:         []string{"key-one"},
	}
	expectedFSlice = []faas.Function{expectedF}

	kName       = "test-key"
	expectedKey = faas.Key{
		Name:        kName,
		Description: "description",
		Functions: []faas.KeysFunction{
			{
				Name:      "function",
				Namespace: "namespace",
			},
		},
		Status: "active",
	}
	expectedKeysSlice  = []faas.Key{expectedKey}
	expectedUpdatedKey = faas.Key{
		Name:        kName,
		Description: "long string",
		Functions: []faas.KeysFunction{
			{
				Name:      "function1",
				Namespace: "namespace1",
			},
			{
				Name:      "function2",
				Namespace: "namespace1",
			},
		},
		Status: "active",
	}
)
