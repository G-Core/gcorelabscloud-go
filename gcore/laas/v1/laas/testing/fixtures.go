package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/laas/v1/laas"

const getStatusResponse = `
{
  "namespace": "2ehwccnfzytsnt576dkmvs",
  "is_initialized": true
}
`

const updateRequest = `
{
  "is_initialized": true
}
`

const getKafkaResponse = `
{
  "hosts": [
    "kafka.example.com",
    "kafka2.example.com",
    "kafka3.example.com"
  ]
}
`

const getOpensearchResponse = `
{
  "hosts": [
    "kafka.example.com",
    "kafka2.example.com",
    "kafka3.example.com"
  ]
}
`

const listTopicResponse = `
{
  "count": 0,
  "results": [
    {
      "name": "string"
    }
  ]
}
`

const createTopicRequest = `
{
  "name": "string"
}
`

const topicResponse = `
{
  "name": "string"
}
`

const regenerateUserResponse = `
{
  "username": "2ehwccnfzytsnt576dkmvs",
  "password": "LnHyPmVor6dAufMtR8GC5WNcNg5NjjAIksjIlFNbaEQ"
}
`

var (
	expectedStatus = laas.Status{
		Namespace:     "2ehwccnfzytsnt576dkmvs",
		IsInitialized: true,
	}
	expectedKafkaHosts      = laas.Hosts{"kafka.example.com", "kafka2.example.com", "kafka3.example.com"}
	expectedOpensearchHosts = laas.Hosts{"kafka.example.com", "kafka2.example.com", "kafka3.example.com"}
	topic                   = laas.Topic{Name: "string"}
	expectedTopicSlice      = []laas.Topic{topic}
	topicName               = "string"

	expectedCreds = laas.User{
		Username: "2ehwccnfzytsnt576dkmvs",
		Password: "LnHyPmVor6dAufMtR8GC5WNcNg5NjjAIksjIlFNbaEQ",
	}
)
