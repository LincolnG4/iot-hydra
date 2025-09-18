package auth

import (
	"errors"
	"strings"
)

type Authenticator interface {
	// returns the method of authentication (sas,token,plain)
	AuthMethod() string

	// validate if all fields are correct. if not correct, it returns an error
	Validate() error
}

/*
* Plain Text
 */

const BasicType = "plain"

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

const TokenType = "token"

type Token struct {
	Token string `json:"token" yaml:"token" validate:"required"`
}

func (t *Token) AuthMethod() string { return TokenType }

func (t *Token) Validate() error {
	// Validate token
	if strings.TrimSpace(t.Token) == "" {
		return errors.New("token cannot be empty")
	}

	return nil
}
