package secrets

import (
	"encoding/json"
	"fmt"
)

type SecretType string

const (
	SymmetricSecretType   SecretType = "symmetric"
	PublicSecretType      SecretType = "public"
	PrivateSecretType     SecretType = "private"
	PassphraseSecretType  SecretType = "passphrase"
	CertificateSecretType SecretType = "certificate"
	OpaqueSecretType      SecretType = "opaque"
)

func (s SecretType) String() string {
	return string(s)
}

func (s SecretType) List() []SecretType {
	return []SecretType{
		SymmetricSecretType,
		PublicSecretType,
		PrivateSecretType,
		PassphraseSecretType,
		CertificateSecretType,
		OpaqueSecretType,
	}
}

func (s SecretType) StringList() []string {
	var sg []string
	for _, v := range s.List() {
		sg = append(sg, v.String())
	}
	return sg
}

func (s SecretType) IsValid() error {
	switch s {
	case SymmetricSecretType,
		PublicSecretType,
		PrivateSecretType,
		PassphraseSecretType,
		CertificateSecretType,
		OpaqueSecretType:
		return nil
	}
	return fmt.Errorf("invalid SecretType type: %v", s)
}

func (s SecretType) ValidOrNil() (*SecretType, error) {
	if s.String() == "" {
		return nil, nil
	}
	err := s.IsValid()
	if err != nil {
		return &s, err
	}
	return &s, nil
}

// UnmarshalJSON - implements Unmarshaler interface for SecretType
func (s *SecretType) UnmarshalJSON(data []byte) error {
	var sg string
	if err := json.Unmarshal(data, &sg); err != nil {
		return err
	}
	v := SecretType(sg)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// MarshalJSON - implements Marshaler interface for SecretType
func (s *SecretType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
