package config

// Config represents the application configuration.
type Config struct {
	CPUThreshold   float64 `json:"cpu_threshold"`
	RAMThreshold   float64 `json:"ram_threshold"`
	DiskThreshold  float64 `json:"disk_threshold"`
	AlertMethod    string  `json:"alert_method"`
	EmailRecipient string  `json:"email_recipient"`
	SlackWebhook   string  `json:"slack_webhook"`
	ServiceName    string  `json:"service_name"`
	Interval       float64 `json:"interval"`
	AlertThreshold float64 `json:"alert_threshold"`
}

var config *Config

// SetConfig stores the parsed configuration.
func SetConfig(cfg *Config) {
	config = cfg
}

// GetConfig retrieves the parsed configuration.
func GetConfig() *Config {
	return config
}
