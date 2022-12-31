package service

type Service struct {
	Config
}

// New returns a new Service instance
func New(options ...Option) *Service {
	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	return &Service{
		Config: cfg,
	}
}

func (s *Service) Start() error {
	return nil
}
