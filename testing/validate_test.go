package testing

import (
	"fmt"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"github.com/stretchr/testify/require"
)

type enum string

func (e enum) IsValid() error {
	if e == "x" {
		return nil
	}
	return fmt.Errorf("invalid Enum type: %v", e)
}

func (e enum) StringList() []string {
	return []string{"x"}
}

func TestValidateEnumTag(t *testing.T) {
	type test struct {
		Value enum `validate:"required,enum"`
	}
	ts1 := test{Value: "y"}
	err := gcorecloud.Validate.Struct(ts1)
	require.Error(t, err)
	err = gcorecloud.TranslateValidationError(err)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not valid")

	ts2 := test{Value: "x"}
	err = gcorecloud.Validate.Struct(ts2)
	require.NoError(t, err)
}
