package auth

type Authenticator interface {
	Apply(opts *ConnectOptions) error
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

type BasicAuth struct {
	Username string
	Password string
}

func (b BasicAuth) Apply(opts *ConnectOptions) error {
	opts.Username = b.Username
	opts.Password = b.Password
	return nil
}

type NatsToken struct {
	Token string
}

func (t NatsToken) Apply(opts *ConnectOptions) error {
	opts.Token = t.Token
	return nil
}
