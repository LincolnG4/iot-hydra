package config

type Service struct {
	Address string `yaml:"address" validate:"required"`
}
