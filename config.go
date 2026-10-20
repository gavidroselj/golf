package golf

type Config struct {
	BasePath string `json:"base_path"`
}

func defaultConfig() Config {
	return Config{
		BasePath: "testdata",
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
