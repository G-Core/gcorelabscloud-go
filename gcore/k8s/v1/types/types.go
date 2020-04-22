package types

import (
	"encoding/json"
	"fmt"
)

type PoolRole string

const (
	NodegroupMasterRole PoolRole = "master"
	NodegroupWorkerRole PoolRole = "worker"
)

func (ng PoolRole) IsValid() error {
	switch ng {
	case NodegroupMasterRole,
		NodegroupWorkerRole:
		return nil
	}
	return fmt.Errorf("invalid PoolRole type: %v", ng)
}

func (ng PoolRole) ValidOrNil() (*PoolRole, error) {
	if ng.String() == "" {
		return nil, nil
	}
	err := ng.IsValid()
	if err != nil {
		return &ng, err
	}
	return &ng, nil
}

func (ng PoolRole) String() string {
	return string(ng)
}

func (ng PoolRole) List() []PoolRole {
	return []PoolRole{
		NodegroupMasterRole,
		NodegroupWorkerRole,
	}
}

func (ng PoolRole) StringList() []string {
	var s []string
	for _, v := range ng.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ng *PoolRole) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := PoolRole(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ng = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (ng *PoolRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(ng.String())
}
