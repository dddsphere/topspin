package config

type (
	Config struct {
		Server
		Mongo
		NATS
		Logging
	}

	Server struct {
		JSONAPIPort int
	}

	Mongo struct {
		Host       string
		Port       int
		User       string
		Pass       string
		Database   string
		MaxRetries int
	}

	NATS struct {
		Host string
		Port int
	}

	Logging struct {
		Level string
	}
)

func (m *Mongo) MaxRetriesUInt64() uint64 {
	return uint64(m.MaxRetries)
}
