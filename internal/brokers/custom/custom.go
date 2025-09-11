package brokers

type CustomConfig struct {
	Endpoint string            `json:"endpoint"`
	Headers  map[string]string `json:"headers"`
	Enabled  bool              `json:"enabled"`
}
