package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LincolnG4/iot-hydra/internal/config"
)

type Authenticator interface {
	// returns the method of authentication (sas,token,plain)
	AuthMethod() string

	// validate if all fields are correct. if not correct, it returns an error
	Validate() error
}

const (
	BasicType = "basic"
	TokenType = "token"
)

// NewAuthenticator acts as a factory for creating an Authenticator.
func NewAuthenticator(cfg config.AuthYAML) (Authenticator, error) {
	switch cfg.Method {
	case BasicType:
		b := &BasicAuth{
			Username: cfg.User,
			Password: cfg.Password,
		}
		if err := b.Validate(); err != nil {
			return nil, err
		}
		return b, nil
	case TokenType:
		t := &TokenAuth{
			Token: cfg.Token,
		}
		if err := t.Validate(); err != nil {
			return nil, err
		}
		return t, nil
	default:
		return nil, fmt.Errorf("authentication method '%s' is not supported", cfg.Method)
	}
}

/*
* Plain Text
 */

type BasicAuth struct {
	Username string `json:"user" yaml:"user" validate:"required_with=Password"`
	Password string `json:"password" yaml:"password" validate:"required_with=Username"`
}

func (b *BasicAuth) AuthMethod() string { return BasicType }

func (b *BasicAuth) Validate() error {
	if strings.TrimSpace(b.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if strings.TrimSpace(b.Password) == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

/*
* Token
 */

type TokenAuth struct {
	Token string `json:"token" yaml:"token" validate:"required"`
}

func (t *TokenAuth) AuthMethod() string { return TokenType }

func (t *TokenAuth) Validate() error {
	// Validate token
	if strings.TrimSpace(t.Token) == "" {
		return errors.New("token cannot be empty")
	}

	return nil
}
