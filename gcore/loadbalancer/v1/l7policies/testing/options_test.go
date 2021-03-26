package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/l7policies"
)

func TestCreateOpts(t *testing.T) {
	options := l7policies.CreateOpts{
		ListenerID: pid,
		Action:     l7policies.ActionReject,
	}
	_, err := options.ToL7PolicyCreateMap()
	require.NoError(t, err)

	options = l7policies.CreateOpts{
		ListenerID: pid,
		Action:     l7policies.ActionRedirectToURL,
	}
	_, err = options.ToL7PolicyCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "RedirectURL")

	options = l7policies.CreateOpts{
		ListenerID:  pid,
		Action:      l7policies.ActionRedirectToURL,
		RedirectURL: "http://www.example.com",
	}
	_, err = options.ToL7PolicyCreateMap()
	require.NoError(t, err)

	options = l7policies.CreateOpts{
		ListenerID: pid,
		Name:       "test",
	}
	_, err = options.ToL7PolicyCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Action")
}

func TestCreateRuleOpts(t *testing.T) {
	options := l7policies.CreateRuleOpts{
		CompareType: l7policies.CompareTypeEqual,
		Value:       "test",
		Type:        l7policies.TypePath,
	}
	_, err := options.ToRuleCreateMap()
	require.NoError(t, err)

	options = l7policies.CreateRuleOpts{
		CompareType: l7policies.CompareTypeEqual,
		Type:        l7policies.TypePath,
	}
	_, err = options.ToRuleCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Value")

	options = l7policies.CreateRuleOpts{
		CompareType: l7policies.CompareTypeEqual,
		Value:       "test",
		Type:        l7policies.RuleType("test"),
	}
	_, err = options.ToRuleCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Type")
}
