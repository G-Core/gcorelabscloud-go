package l7policies

import (
	"encoding/json"
	"fmt"
)

type Action string
type RuleType string
type CompareType string

const (
	ActionRedirectToPool Action = "REDIRECT_TO_POOL"
	ActionRedirectToURL  Action = "REDIRECT_TO_URL"
	ActionRedirectPrefix Action = "REDIRECT_PREFIX"
	ActionReject         Action = "REJECT"

	TypeCookie          RuleType = "COOKIE"
	TypeFileType        RuleType = "FILE_TYPE"
	TypeHeader          RuleType = "HEADER"
	TypeHostName        RuleType = "HOST_NAME"
	TypePath            RuleType = "PATH"
	TypeSSLConnHasCert  RuleType = "SSL_CONN_HAS_CERT"
	TypeSSLVerifyResult RuleType = "SSL_VERIFY_RESULT"
	TypeSSLDNField      RuleType = "SSL_DN_FIELD"

	CompareTypeContains  CompareType = "CONTAINS"
	CompareTypeEndWith   CompareType = "ENDS_WITH"
	CompareTypeEqual     CompareType = "EQUAL_TO"
	CompareTypeRegex     CompareType = "REGEX"
	CompareTypeStartWith CompareType = "STARTS_WITH"
)

func (a Action) IsValid() error {
	switch a {
	case ActionRedirectToPool,
		ActionRedirectToURL,
		ActionReject,
		ActionRedirectPrefix:
		return nil
	}
	return fmt.Errorf("invalid Action: %v", a)
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
		ActionRedirectToPool,
		ActionRedirectToURL,
		ActionReject,
		ActionRedirectPrefix,
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

func (rt RuleType) IsValid() error {
	switch rt {
	case TypeCookie,
		TypeFileType,
		TypeHeader,
		TypeHostName,
		TypePath,
		TypeSSLConnHasCert,
		TypeSSLVerifyResult,
		TypeSSLDNField:
		return nil
	}
	return fmt.Errorf("invalid rule type: %v", rt)
}

func (rt RuleType) ValidOrNil() (*RuleType, error) {
	if rt.String() == "" {
		return nil, nil
	}
	err := rt.IsValid()
	if err != nil {
		return &rt, err
	}
	return &rt, nil
}

func (rt RuleType) String() string {
	return string(rt)
}

func (rt RuleType) List() []RuleType {
	return []RuleType{
		TypeCookie,
		TypeFileType,
		TypeHeader,
		TypeHostName,
		TypePath,
		TypeSSLConnHasCert,
		TypeSSLVerifyResult,
		TypeSSLDNField,
	}
}

func (rt RuleType) StringList() []string {
	var s []string
	for _, v := range rt.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (rt *RuleType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := RuleType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*rt = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (rt *RuleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(rt.String())
}

func (ct CompareType) IsValid() error {
	switch ct {
	case CompareTypeContains,
		CompareTypeEndWith,
		CompareTypeEqual,
		CompareTypeRegex,
		CompareTypeStartWith:
		return nil
	}
	return fmt.Errorf("invalid compare type: %v", ct)
}

func (ct CompareType) ValidOrNil() (*CompareType, error) {
	if ct.String() == "" {
		return nil, nil
	}
	err := ct.IsValid()
	if err != nil {
		return &ct, err
	}
	return &ct, nil
}

func (ct CompareType) String() string {
	return string(ct)
}

func (ct CompareType) List() []CompareType {
	return []CompareType{
		CompareTypeContains,
		CompareTypeEndWith,
		CompareTypeEqual,
		CompareTypeRegex,
		CompareTypeStartWith,
	}
}

func (ct CompareType) StringList() []string {
	var s []string
	for _, v := range ct.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ct *CompareType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := CompareType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ct = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (ct *CompareType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}
