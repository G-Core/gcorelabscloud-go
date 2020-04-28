package testing

import (
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcoreclient/utils"
	"github.com/stretchr/testify/require"
)

func TestValidateEqualSlicesLength(t *testing.T) {
	strSlice := []string{
		"one",
		"two",
	}
	intSlice := []int{
		1,
		3,
		3,
	}
	err := utils.ValidateEqualSlicesLength(strSlice, intSlice)
	require.Error(t, err)
	strSlice = []string{
		"one",
		"two",
	}
	intSlice = []int{
		1,
		3,
	}
	err = utils.ValidateEqualSlicesLength(strSlice, intSlice)
	require.NoError(t, err)
}
