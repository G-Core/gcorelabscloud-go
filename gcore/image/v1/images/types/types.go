package types

import (
	"encoding/json"
	"fmt"
)

type Visibility string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityShared  Visibility = "shared"
	VisibilityPublic  Visibility = "public"
)

func (v Visibility) IsValid() error {
	switch v {
	case VisibilityPrivate, VisibilityShared, VisibilityPublic:
		return nil
	}
	return fmt.Errorf("invalid Visibility type: %v", v)
}

func (v Visibility) ValidOrNil() (*Visibility, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v Visibility) String() string {
	return string(v)
}

func (v Visibility) List() []Visibility {
	return []Visibility{VisibilityPrivate, VisibilityShared, VisibilityPublic}
}

func (v Visibility) StringList() []string {
	var s []string
	for _, v := range v.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for VolumeSource
func (v *Visibility) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := Visibility(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for VolumeSource
func (v *Visibility) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}
