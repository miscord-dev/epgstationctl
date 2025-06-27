package config

type Config struct {
	Server ServerConfig
	Output OutputConfig
	Debug  bool
}

type ServerConfig struct {
	URL     string
	Timeout int
}

type OutputConfig struct {
	Format   string
	NoHeader bool
}

func NewConfig(serverURL string, timeout int, outputFormat string, noHeader bool, debug bool) *Config {
	// Set defaults if not provided
	if serverURL == "" {
		serverURL = "http://localhost:8888"
	}
	if timeout <= 0 {
		timeout = 30
	}
	if outputFormat == "" {
		outputFormat = "table"
	}

	return &Config{
		Server: ServerConfig{
			URL:     serverURL,
			Timeout: timeout,
		},
		Output: OutputConfig{
			Format:   outputFormat,
			NoHeader: noHeader,
		},
		Debug: debug,
	}
}

func (c *Config) GetAPIBaseURL() string {
	return c.Server.URL + "/api"
}
