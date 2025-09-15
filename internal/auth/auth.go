package auth

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
	Username string
	Password string
}

func (b BasicAuth) Apply(opts *ConnectOptions) error {
	opts.Username = b.Username
	opts.Password = b.Password
	return nil
}

func (b BasicAuth) Type() string {
	return BasicType
}

const NatsTokenType = "natsToken"

type NatsToken struct {
	Token string
}

func (n NatsToken) Apply(opts *ConnectOptions) error {
	opts.Token = n.Token
	return nil
}

func (n NatsToken) Type() string {
	return NatsTokenType
}
