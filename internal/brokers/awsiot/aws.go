package brokers

type AWSConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	ThingName string `json:"thing_name"`
	Enabled   bool   `json:"enabled"`
}
