package config

type DBConfig struct {
	Dialect string
	URL     string
}

type APIConfig struct {
	Host string
	Port int
}

type Config struct {
	DBConfig
	APIConfig
}
