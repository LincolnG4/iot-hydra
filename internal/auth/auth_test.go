package auth

import (
	"reflect"
	"strings"
	"testing"
)

func TestBasicAuth_AuthMethod(t *testing.T) {
	b := BasicAuth{}
	expected := BasicType
	if b.AuthMethod() != expected {
		t.Errorf("BasicAuth.AuthMethod() got = %s, want %s", b.AuthMethod(), expected)
	}
}

func TestBasicAuth_validate(t *testing.T) {
	tests := []struct {
		name     string
		auth     BasicAuth
		wantErr  bool
		errMatch string
	}{
		{
			name:    "Valid credentials",
			auth:    BasicAuth{Username: "testuser", Password: "testpassword"},
			wantErr: false,
		},
		{
			name:     "Empty username",
			auth:     BasicAuth{Username: "", Password: "testpassword"},
			wantErr:  true,
			errMatch: "username cannot be empty",
		},
		{
			name:     "Whitespace username",
			auth:     BasicAuth{Username: "   ", Password: "testpassword"},
			wantErr:  true,
			errMatch: "username cannot be empty",
		},
		{
			name:     "Empty password",
			auth:     BasicAuth{Username: "testuser", Password: ""},
			wantErr:  true,
			errMatch: "password cannot be empty",
		},
		{
			name:     "Whitespace password",
			auth:     BasicAuth{Username: "testuser", Password: "   "},
			wantErr:  true,
			errMatch: "password cannot be empty",
		},
		{
			name:     "Both empty",
			auth:     BasicAuth{Username: "", Password: ""},
			wantErr:  true,
			errMatch: "username cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.auth.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicAuth.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMatch != "" && !strings.Contains(err.Error(), tt.errMatch) {
				t.Errorf("BasicAuth.validate() error = %v, want error containing %q", err, tt.errMatch)
			}
		})
	}
}

func TestToken_AuthMethod(t *testing.T) {
	tok := TokenAuth{}
	expected := TokenType
	if tok.AuthMethod() != expected {
		t.Errorf("Token.AuthMethod() got = %s, want %s", tok.AuthMethod(), expected)
	}
}

func TestToken_validate(t *testing.T) {
	tests := []struct {
		name     string
		auth     TokenAuth
		wantErr  bool
		errMatch string
	}{
		{
			name:    "Valid token",
			auth:    TokenAuth{Token: "some-secret-token"},
			wantErr: false,
		},
		{
			name:     "Empty token",
			auth:     TokenAuth{Token: ""},
			wantErr:  true,
			errMatch: "token cannot be empty",
		},
		{
			name:     "Whitespace token",
			auth:     TokenAuth{Token: "   "},
			wantErr:  true,
			errMatch: "token cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.auth.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMatch != "" && !strings.Contains(err.Error(), tt.errMatch) {
				t.Errorf("Token.validate() error = %v, want error containing %q", err, tt.errMatch)
			}
		})
	}
}

func TestAuthenticatorInterface(t *testing.T) {
	auths := []Authenticator{
		&BasicAuth{},
		&TokenAuth{},
		// Add other Authenticator implementations here as they are created
	}

	for _, a := range auths {
		t.Run(reflect.TypeOf(a).String(), func(t *testing.T) {
			// Test that AuthMethod doesn't panic and returns a non-empty string
			method := a.AuthMethod()
			if method == "" {
				t.Errorf("%T.AuthMethod() returned an empty string", a)
			}

			_ = a.Validate()
		})
	}
}
