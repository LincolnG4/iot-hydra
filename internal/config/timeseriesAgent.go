package config

type TelemetryAgentYAML struct {
	QueueSize int          `yaml:"queueSize" validate:"gt=0"`
	Brokers   []BrokerYAML `yaml:"brokers" validate:"required,min=1,dive"`
}

type BrokerYAML struct {
	Name    string   `yaml:"name" validate:"required"`
	Type    string   `yaml:"type" validate:"required"`
	Address string   `yaml:"address" validate:"required"`
	Auth    AuthYAML `yaml:"auth" validate:"required"`
}
type AuthYAML struct {
	Method   string `yaml:"method" validate:"required"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
}
