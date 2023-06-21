package servergroups

import (
	"encoding/json"
	"fmt"
)

type ServerGroupPolicy string

const (
	AffinityPolicy     ServerGroupPolicy = "affinity"
	AntiAffinityPolicy ServerGroupPolicy = "anti-affinity"
	SoftAffinityPolicy ServerGroupPolicy = "soft-anti-affinity"
)

func (s ServerGroupPolicy) String() string {
	return string(s)
}

func (s ServerGroupPolicy) List() []ServerGroupPolicy {
	return []ServerGroupPolicy{AffinityPolicy, AntiAffinityPolicy, SoftAffinityPolicy}
}

func (s ServerGroupPolicy) StringList() []string {
	var sg []string
	for _, v := range s.List() {
		sg = append(sg, v.String())
	}
	return sg
}

func (s ServerGroupPolicy) IsValid() error {
	switch s {
	case AffinityPolicy, AntiAffinityPolicy, SoftAffinityPolicy:
		return nil
	}
	return fmt.Errorf("invalid ServerGroupPolicy type: %v", s)
}

func (s ServerGroupPolicy) ValidOrNil() (*ServerGroupPolicy, error) {
	if s.String() == "" {
		return nil, nil
	}
	err := s.IsValid()
	if err != nil {
		return &s, err
	}
	return &s, nil
}

// UnmarshalJSON - implements Unmarshaler interface for ServerGroupPolicy
func (s *ServerGroupPolicy) UnmarshalJSON(data []byte) error {
	var sg string
	if err := json.Unmarshal(data, &sg); err != nil {
		return err
	}
	v := ServerGroupPolicy(sg)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// MarshalJSON - implements Marshaler interface for ServerGroupPolicy
func (s *ServerGroupPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
