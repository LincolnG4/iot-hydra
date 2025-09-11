package brokers

type AzureConfig struct {
	ConnectionString string `json:"connection_string"`
	DeviceID         string `json:"device_id"`
	Enabled          bool   `json:"enabled"`
}
