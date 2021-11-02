package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type RoleIDType int
type RoleNameType string

const (
	RoleIDAdministrators RoleIDType = 1
	RoleIDUsers          RoleIDType = 2
	RoleIDEngineers      RoleIDType = 5
	RoleIDAPIWeb         RoleIDType = 3009
	RoleIDAPI            RoleIDType = 3022

	RoleNameUsers          RoleNameType = "Users"
	RoleNameAdministrators RoleNameType = "Administrators"
	RoleNameEngineers      RoleNameType = "Engineers"
	RoleNameAPI            RoleNameType = "Purge and Prefetch only (API)"
	RoleNameAPIWeb         RoleNameType = "Purge and Prefetch only (API+Web)"
)

func (r RoleIDType) IsValid() error {
	switch r {
	case RoleIDAPI, RoleIDAPIWeb, RoleIDAdministrators, RoleIDUsers, RoleIDEngineers:
		return nil
	}
	return fmt.Errorf("invalid RoleIDType type: %v", r)
}

func (r RoleIDType) ValidOrNil() (*RoleIDType, error) {
	if r.String() == "" {
		return nil, nil
	}
	err := r.IsValid()
	if err != nil {
		return &r, err
	}
	return &r, nil
}

func (r RoleIDType) String() string {
	return strconv.Itoa(int(r))
}

func (r RoleIDType) List() []RoleIDType {
	return []RoleIDType{RoleIDAPI, RoleIDAPIWeb, RoleIDAdministrators, RoleIDUsers, RoleIDEngineers}
}

func (r RoleIDType) StringList() []string {
	var s []string
	for _, v := range r.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for RoleIDType
func (r *RoleIDType) UnmarshalJSON(data []byte) error {
	var s int
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := RoleIDType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*r = v
	return nil
}

// MarshalJSON - implements Marshaler interface for RoleIDType
func (r *RoleIDType) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*r))
}

func (rs RoleNameType) IsValid() error {
	switch rs {
	case RoleNameAdministrators,
		RoleNameUsers,
		RoleNameEngineers,
		RoleNameAPI,
		RoleNameAPIWeb:
		return nil
	}
	return fmt.Errorf("invalid RoleNameType type: %v", rs)
}

func (rs RoleNameType) ValidOrNil() (*RoleNameType, error) {
	if rs.String() == "" {
		return nil, nil
	}
	err := rs.IsValid()
	if err != nil {
		return &rs, err
	}
	return &rs, nil
}

func (rs RoleNameType) String() string {
	return string(rs)
}

func (rs RoleNameType) List() []RoleNameType {
	return []RoleNameType{
		RoleNameAdministrators,
		RoleNameUsers,
		RoleNameEngineers,
		RoleNameAPI,
		RoleNameAPIWeb,
	}
}

func (rs RoleNameType) StringList() []string {
	var s []string
	for _, v := range rs.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for RoleNameType
func (rs *RoleNameType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := RoleNameType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*rs = v
	return nil
}

// MarshalJSON - implements Marshaler interface for RoleNameType
func (rs *RoleNameType) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}
