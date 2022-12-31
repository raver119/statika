package service

type Option func(cfg *Config)

// WithPredefinedConfig returns an Option that sets the predefined Config
func WithPredefinedConfig(cfg *Config) Option {
	return func(c *Config) {
		*c = *cfg
	}
}

func WithPutMiddleware() {

}

func WithGetMiddleware() {

}
