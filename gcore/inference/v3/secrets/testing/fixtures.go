package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/secrets"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "name": "aws-dev",
      "type": "aws-iam",
      "data": {
		"aws_access_key_id": "aws_access_key_id",
		"aws_secret_access_key": "aws_secret_access_key"
      }
    }
  ]
}
`

const CreateRequest = `
    {
      "name": "aws-dev",
      "type": "aws-iam",
      "data": {
		"aws_access_key_id": "aws_access_key_id",
		"aws_secret_access_key": "aws_secret_access_key"
      }
    }
`

const UpdateRequest = `
    {
      "type": "aws-iam",
      "data": {
		"aws_access_key_id": "aws_access_key_id",
		"aws_secret_access_key": "aws_secret_access_key"
      }
    }
`

const GetResponse = `
    {
      "name": "aws-dev",
      "type": "aws-iam",
      "data": {
		"aws_access_key_id": "aws_access_key_id",
		"aws_secret_access_key": "aws_secret_access_key"
      }
    }
`

var (
	Secret1 = secrets.InferenceSecret{
		Name: "aws-dev",
		Type: "aws-iam",
		Data: secrets.SecretData{
			AWSSecretKeyID:     "aws_access_key_id",
			AWSSecretAccessKey: "aws_secret_access_key",
		},
	}
	SecretSlice = []secrets.InferenceSecret{Secret1}
)
