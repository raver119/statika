package service

type Config struct {
	Backends      []string
	Encryption    bool
	ListenAddress string
}

var DefaultConfig = Config{
	Backends:      nil,
	Encryption:    true,
	ListenAddress: "127.0.0.1:11911",
}
