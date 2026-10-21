package golf

type Config struct {
	BasePath                string `json:"base_path"`
	DynamicValuePlaceholder string `json:"dynamic_value_placeholder"`
}

func defaultConfig() Config {
	return Config{
		BasePath:                "testdata",
		DynamicValuePlaceholder: "-- DYNAMIC VALUE --",
	}
}

var config = defaultConfig()

type Option func(*Config)

func Configure(opts ...Option) {
	for _, opt := range opts {
		opt(&config)
	}
}

func WithBasePath(basePath string) Option {
	return func(c *Config) {
		c.BasePath = basePath
	}
}

func WithDynamicValuePlaceholder(dynamicValuePlaceholder string) Option {
	return func(c *Config) {
		c.DynamicValuePlaceholder = dynamicValuePlaceholder
	}
}
