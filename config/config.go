package config

type Config struct {
	CPU      float32 `json:"cpu"`
	RAM      float32 `json:"ram"`
	Disk     float32 `json:"disk"`
	LogFile  string  `json:"log_file"`
	Interval int     `json:"interval"`
}

// type UserConfig config
