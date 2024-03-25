package types

import (
	"encoding/json"
	"fmt"
)

type RuleDirection string
type EtherType string
type Protocol string
type Action string

const (
	RuleDirectionIngress RuleDirection = "ingress"
	RuleDirectionEgress  RuleDirection = "egress"
	EtherTypeIPv4        EtherType     = "IPv4"
	EtherTypeIPv6        EtherType     = "IPv6"
	ProtocolTCP          Protocol      = "tcp"
	ProtocolUDP          Protocol      = "udp"
	ProtocolICMP         Protocol      = "icmp"
	ProtocolIPv6ICMP     Protocol      = "ipv6-icmp"
	ProtocolIPv6Route    Protocol      = "ipv6-route"
	ProtocolIPv6Opts     Protocol      = "ipv6-opts"
	ProtocolIPv6Nonxt    Protocol      = "ipv6-nonxt"
	ProtocolIPv6Frag     Protocol      = "ipv6-frag"
	ProtocolIPv6Encap    Protocol      = "ipv6-encap"
	ProtocolAny          Protocol      = "any"
	ProtocolAH           Protocol      = "ah"
	ProtocolDCCP         Protocol      = "dccp"
	ProtocolEGP          Protocol      = "egp"
	ProtocolESP          Protocol      = "esp"
	ProtocolGRE          Protocol      = "gre"
	ProtocolIGMP         Protocol      = "igmp"
	ProtocolOSPF         Protocol      = "ospf"
	ProtocolPGM          Protocol      = "pgm"
	ProtocolRSVP         Protocol      = "rsvp"
	ProtocolSCTP         Protocol      = "sctp"
	ProtocolUDPLITE      Protocol      = "udplite"
	ProtocolVRRP         Protocol      = "vrrp"
	Protocol51           Protocol      = "51"
	Protocol50           Protocol      = "50"
	Protocol112          Protocol      = "112"
	Protocol0            Protocol      = "0"
	ProtocolIPinIP       Protocol      = "4"
	ProtocolIPIP         Protocol      = "ipip"
	ProtocolIPEncap      Protocol      = "ipencap"
	ActionCreate         Action        = "create"
	ActionDelete         Action        = "delete"
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
		ProtocolIPv6ICMP,
		ProtocolIPv6Route,
		ProtocolIPv6Opts,
		ProtocolIPv6Nonxt,
		ProtocolIPv6Frag,
		ProtocolIPv6Encap,
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
		ProtocolVRRP,
		Protocol51,
		Protocol50,
		Protocol0,
		Protocol112,
		ProtocolIPinIP,
		ProtocolIPEncap,
		ProtocolIPIP:
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
		ProtocolIPv6ICMP,
		ProtocolIPv6Route,
		ProtocolIPv6Opts,
		ProtocolIPv6Nonxt,
		ProtocolIPv6Frag,
		ProtocolIPv6Encap,
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
		Protocol51,
		Protocol50,
		Protocol112,
		Protocol0,
		ProtocolIPinIP,
		ProtocolIPIP,
		ProtocolIPEncap,
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

func (a Action) IsValid() error {
	switch a {
	case ActionCreate,
		ActionDelete:
		return nil
	}
	return fmt.Errorf("invalid Action type: %v", a)
}

func (a Action) ValidOrNil() (*Action, error) {
	if a.String() == "" {
		return nil, nil
	}
	err := a.IsValid()
	if err != nil {
		return &a, err
	}
	return &a, nil
}

func (a Action) String() string {
	return string(a)
}

func (a Action) List() []Action {
	return []Action{
		ActionCreate,
		ActionDelete,
	}
}

func (a Action) StringList() []string {
	var s []string
	for _, v := range a.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (a *Action) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := Action(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*a = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (a *Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}
