package value_object

import (
	"errors"
	"regexp"
)

var ErrInvalidSshKey = errors.New("ssh key is invalid")

const (
	SshRegexp string = `^(ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519|ssh-rsa)\s+(?:[a-zA-Z0-9+/]{4})*(?:|[a-zA-Z0-9+/]{3}=|[a-zA-Z0-9+/]{2}==|[a-zA-Z0-9+/]===)[\s+\x21-\x7F]+$`
)

type SshKey struct {
	value  string
	Regexp string
}

func NewSshKey(value string) (*SshKey, error) {
	r, _ := regexp.Compile(SshRegexp)
	if !r.MatchString(value) {
		return nil, ErrInvalidSshKey
	}

	return &SshKey{value: value}, nil
}

func (s SshKey) String() string {
	return s.value
}
