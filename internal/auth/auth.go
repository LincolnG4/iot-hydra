package auth

import (
	"errors"
	"strings"
)

type Authenticator interface {
	Apply(opts *ConnectOptions) error
	Type() string
}

// ConnectOptions is a generic bag for connection parameters.
// Each broker adapter can interpret what it needs.
type ConnectOptions struct {
	Username string
	Password string
	Token    string
	CertFile string
	KeyFile  string
}

const BasicType = "plain"

type BasicAuth struct {
	Username string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}

func (b BasicAuth) Apply(opts *ConnectOptions) error {
	// Validate required fields
	if strings.TrimSpace(b.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if strings.TrimSpace(b.Password) == "" {
		return errors.New("password cannot be empty")
	}
	opts.Username = strings.TrimSpace(b.Username)
	opts.Password = b.Password

	return opts.Validate()
}

func (b BasicAuth) Type() string {
	return BasicType
}

const NatsTokenType = "natsToken"

type NatsToken struct {
	Token string `json:"token" yaml:"token"`
}

func (n NatsToken) Apply(opts *ConnectOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	// Validate token
	if strings.TrimSpace(n.Token) == "" {
		return errors.New("token cannot be empty")
	}

	opts.Token = strings.TrimSpace(n.Token)
	return opts.Validate()
}

func (n NatsToken) Type() string {
	return NatsTokenType
}

func (opts *ConnectOptions) Validate() error {
	// Check for conflicting auth methods
	authMethods := 0
	if opts.Username != "" || opts.Password != "" {
		authMethods++
	}
	if opts.Token != "" {
		authMethods++
	}
	if opts.CertFile != "" || opts.KeyFile != "" {
		authMethods++
	}

	if authMethods == 0 {
		return errors.New("no authentication method provided")
	}
	if authMethods > 1 {
		return errors.New("multiple authentication methods provided, only one is allowed")
	}

	// Validate cert/key pair completeness
	if (opts.CertFile != "" && opts.KeyFile == "") || (opts.CertFile == "" && opts.KeyFile != "") {
		return errors.New("both CertFile and KeyFile must be provided for certificate authentication")
	}

	return nil
}
