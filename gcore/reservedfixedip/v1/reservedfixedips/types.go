package reservedfixedips

import (
	"encoding/json"
	"fmt"
)

type ReservedFixedIPType string
type IPFamilyType string

const (
	External  ReservedFixedIPType = "external"
	Subnet    ReservedFixedIPType = "subnet"
	AnySubnet ReservedFixedIPType = "any_subnet"
	IPAddress ReservedFixedIPType = "ip_address"

	IPv4IPFamilyType      IPFamilyType = "ipv4"
	IPv6IPFamilyType      IPFamilyType = "ipv6"
	DualStackIPFamilyType IPFamilyType = "dual"
)

func (t ReservedFixedIPType) String() string {
	return string(t)
}

func (t ReservedFixedIPType) List() []ReservedFixedIPType {
	return []ReservedFixedIPType{External, Subnet, AnySubnet, IPAddress}
}

func (t ReservedFixedIPType) StringList() []string {
	var s []string
	for _, v := range t.List() {
		s = append(s, v.String())
	}
	return s
}

func (t ReservedFixedIPType) IsValid() error {
	switch t {
	case External, Subnet, AnySubnet, IPAddress:
		return nil
	}
	return fmt.Errorf("invalid ReservedFixedIPType type: %v", t)
}

func (t ReservedFixedIPType) ValidOrNil() (*ReservedFixedIPType, error) {
	if t.String() == "" {
		return nil, nil
	}
	err := t.IsValid()
	if err != nil {
		return &t, err
	}
	return &t, nil
}

// UnmarshalJSON - implements Unmarshaler interface for ReservedFixedIPType
func (t *ReservedFixedIPType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := ReservedFixedIPType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*t = v
	return nil
}

// MarshalJSON - implements Marshaler interface for ReservedFixedIPType
func (t *ReservedFixedIPType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
