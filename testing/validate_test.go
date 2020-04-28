package testing

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go"
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

func TestAllowedWithoutTag(t *testing.T) {
	type test struct {
		Value1 string `validate:"allowed_without=Value2"`
		Value2 string `validate:"allowed_without=Value1"`
	}
	ts1 := test{Value1: "y", Value2: "y"}
	err := gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(ts1))
	require.Error(t, err)
	require.Contains(t, err.Error(), "should not be")
	ts2 := test{Value1: "y", Value2: ""}
	err = gcorecloud.Validate.Struct(ts2)
	require.NoError(t, err)
	ts3 := test{Value1: "", Value2: "y"}
	err = gcorecloud.Validate.Struct(ts3)
	require.NoError(t, err)
	ts4 := test{Value1: "", Value2: ""}
	err = gcorecloud.Validate.Struct(ts4)
	require.NoError(t, err)
}

func TestAllowedWithoutAllTag(t *testing.T) {
	type test struct {
		Value1 string `validate:"allowed_without_all=Value2 Value3"`
		Value2 string `validate:"allowed_without_all=Value1 Value3"`
		Value3 string `validate:"allowed_without_all=Value1 Value2"`
		Value4 string
	}
	ts1 := test{Value1: "y", Value2: "y"}
	err := gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(ts1))
	require.Error(t, err)
	require.Contains(t, err.Error(), "should not be")
	ts2 := test{Value1: "y", Value2: ""}
	err = gcorecloud.Validate.Struct(ts2)
	require.NoError(t, err)
	ts3 := test{Value1: "", Value2: "y"}
	err = gcorecloud.Validate.Struct(ts3)
	require.NoError(t, err)
	ts4 := test{Value1: "", Value2: ""}
	err = gcorecloud.Validate.Struct(ts4)
	require.NoError(t, err)
	ts5 := test{Value1: "y", Value4: "y"}
	err = gcorecloud.Validate.Struct(ts5)
	require.NoError(t, err)
	ts6 := test{Value1: "y", Value4: "y", Value2: "y"}
	err = gcorecloud.Validate.Struct(ts6)
	require.Error(t, err)
}
