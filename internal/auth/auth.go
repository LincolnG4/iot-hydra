package auth

import (
	"errors"
	"strings"
)

type Authenticator interface {
	// returns the method of authentication (sas,token,plain)
	AuthMethod() string

	// validate if all fields are correct. if not correct, it returns an error
	validate() error
}

/*
* Plain Text
 */

const BasicType = "plain"

type BasicAuth struct {
	Username string `json:"user" yaml:"user" validate:"required_with=Password"`
	Password string `json:"password" yaml:"password" validate:"required_with=Username"`
}

func (b BasicAuth) AuthMethod() string { return BasicType }

func (b BasicAuth) validate() error {
	// Validate required fields
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

const TokenType = "Token"

type Token struct {
	Token string `json:"token" yaml:"token" validate:"required"`
}

func (n Token) AuthMethod() string { return TokenType }

func (n Token) validate() error {
	// Validate token
	if strings.TrimSpace(n.Token) == "" {
		return errors.New("token cannot be empty")
	}

	return nil
}
