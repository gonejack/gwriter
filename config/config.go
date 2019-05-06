package config

// Config is the config struct use for writer configuration
type Config struct {
	PathTpl      string
	BaseExt      string
	WriteExt     string
	PathInfo     map[string]string
	UpdateMoment string
	UpdatePeriod int
	UpdateSize   int64
}
