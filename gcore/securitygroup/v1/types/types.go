package types

import (
	"encoding/json"
	"fmt"
)

type RuleDirection string
type EtherType string
type Protocol string

const (
	RuleDirectionIngress RuleDirection = "ingress"
	RuleDirectionEgress  RuleDirection = "egress"
	EtherTypeIPv4        EtherType     = "IPv4"
	EtherTypeIPv6        EtherType     = "IPv6"
	ProtocolTCP          Protocol      = "tcp"
	ProtocolUDP          Protocol      = "udp"
	ProtocolICMP         Protocol      = "icmp"
	ProtocolAny          Protocol      = "any"
	ProtocolAH           Protocol      = "ah"
	ProtocolDCCP         Protocol      = "dccp"
	ProtocolEGP          Protocol      = "egp"
	ProtocolESP          Protocol      = "esp"
	ProtocolGRE          Protocol      = "gre"
	ProtocolIGMP         Protocol      = "imgp"
	ProtocolOSPF         Protocol      = "ospf"
	ProtocolPGM          Protocol      = "pgm"
	ProtocolRSVP         Protocol      = "rsvp"
	ProtocolSCTP         Protocol      = "sctp"
	ProtocolUDPLITE      Protocol      = "udplite"
	ProtocolVRRP         Protocol      = "vrrp"
)

func (rd RuleDirection) IsValid() error {
	switch rd {
	case RuleDirectionEgress,
		RuleDirectionIngress:
		return nil
	}
	return fmt.Errorf("invalid RuleDirection type: %v", rd)
}

func (rd RuleDirection) ValidOrNil() (*RuleDirection, error) {
	if rd.String() == "" {
		return nil, nil
	}
	err := rd.IsValid()
	if err != nil {
		return &rd, err
	}
	return &rd, nil
}

func (rd RuleDirection) String() string {
	return string(rd)
}

func (rd RuleDirection) List() []RuleDirection {
	return []RuleDirection{
		RuleDirectionIngress,
		RuleDirectionEgress,
	}
}

func (rd RuleDirection) StringList() []string {
	var s []string
	for _, v := range rd.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (rd *RuleDirection) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := RuleDirection(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*rd = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (rd *RuleDirection) MarshalJSON() ([]byte, error) {
	return json.Marshal(rd.String())
}

func (et EtherType) IsValid() error {
	switch et {
	case EtherTypeIPv4,
		EtherTypeIPv6:
		return nil
	}
	return fmt.Errorf("invalid EtherType: %v", et)
}

func (et EtherType) ValidOrNil() (*EtherType, error) {
	if et.String() == "" {
		return nil, nil
	}
	err := et.IsValid()
	if err != nil {
		return &et, err
	}
	return &et, nil
}

func (et EtherType) String() string {
	return string(et)
}

func (et EtherType) List() []EtherType {
	return []EtherType{
		EtherTypeIPv4,
		EtherTypeIPv6,
	}
}

func (et EtherType) StringList() []string {
	var s []string
	for _, v := range et.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (et *EtherType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := EtherType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*et = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (et *EtherType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

func (p Protocol) IsValid() error {
	switch p {
	case ProtocolICMP,
		ProtocolAny,
		ProtocolTCP,
		ProtocolUDP,
		ProtocolAH,
		ProtocolDCCP,
		ProtocolEGP,
		ProtocolESP,
		ProtocolGRE,
		ProtocolIGMP,
		ProtocolOSPF,
		ProtocolPGM,
		ProtocolRSVP,
		ProtocolSCTP,
		ProtocolUDPLITE,
		ProtocolVRRP:
		return nil
	}
	return fmt.Errorf("invalid Protocol: %v", p)
}

func (p Protocol) ValidOrNil() (*Protocol, error) {
	if p.String() == "" {
		return nil, nil
	}
	err := p.IsValid()
	if err != nil {
		return &p, err
	}
	return &p, nil
}

func (p Protocol) String() string {
	return string(p)
}

func (p Protocol) List() []Protocol {
	return []Protocol{
		ProtocolUDP,
		ProtocolTCP,
		ProtocolAny,
		ProtocolICMP,
		ProtocolAH,
		ProtocolDCCP,
		ProtocolEGP,
		ProtocolESP,
		ProtocolGRE,
		ProtocolIGMP,
		ProtocolOSPF,
		ProtocolPGM,
		ProtocolRSVP,
		ProtocolSCTP,
		ProtocolUDPLITE,
		ProtocolVRRP,
	}
}

func (p Protocol) StringList() []string {
	var s []string
	for _, v := range p.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (p *Protocol) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := Protocol(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*p = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (p *Protocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
