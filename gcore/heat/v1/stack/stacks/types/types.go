package types

import (
	"encoding/json"
	"fmt"
)

// SortDir is a type for specifying in which direction to sort a list of stacks.
type SortDir string

// SortKey is a type for specifying by which key to sort a list of stacks.
type SortKey string

const (
	SortAsc       SortDir = "asc"
	SortDesc      SortDir = "desc"
	SortName      SortKey = "name"
	SortStatus    SortKey = "status"
	SortCreatedAt SortKey = "created_at"
	SortUpdatedAt SortKey = "updated_at"
)

func (sd SortDir) IsValid() error {
	switch sd {
	case SortAsc, SortDesc:
		return nil
	}
	return fmt.Errorf("invalid SortDir type: %v", sd)
}

func (sd SortDir) ValidOrNil() (*SortDir, error) {
	if sd.String() == "" {
		return nil, nil
	}
	err := sd.IsValid()
	if err != nil {
		return &sd, err
	}
	return &sd, nil
}

func (sd SortDir) String() string {
	return string(sd)
}

func (sd SortDir) List() []SortDir {
	return []SortDir{
		SortDesc, SortAsc,
	}
}

func (sd SortDir) StringList() []string {
	var s []string
	for _, v := range sd.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (sd *SortDir) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := SortDir(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*sd = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (sd *SortDir) MarshalJSON() ([]byte, error) {
	return json.Marshal(sd.String())
}

func (sk SortKey) IsValid() error {
	switch sk {
	case SortName, SortStatus, SortCreatedAt, SortUpdatedAt:
		return nil
	}
	return fmt.Errorf("invalid SortKey type: %v", sk)
}

func (sk SortKey) ValidOrNil() (*SortKey, error) {
	if sk.String() == "" {
		return nil, nil
	}
	err := sk.IsValid()
	if err != nil {
		return &sk, err
	}
	return &sk, nil
}

func (sk SortKey) String() string {
	return string(sk)
}

func (sk SortKey) List() []SortKey {
	return []SortKey{
		SortName, SortStatus, SortCreatedAt, SortUpdatedAt,
	}
}

func (sk SortKey) StringList() []string {
	var s []string
	for _, v := range sk.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (sk *SortKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := SortKey(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*sk = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (sk *SortKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(sk.String())
}
