package config

// DBConfig holds the database configuration values, like the database dialect and connection URL or DSN.
type DBConfig struct {
	// Dialect is the database driver dialect (database type).
	Dialect string

	// URL is the connection URL or DSN.
	URL string
}

// APIConfig holds configuration values for the API routing and setup.
type APIConfig struct {
	// Host to bind to. Default is an empty string.
	Host string

	// Port to listen on.
	Port int
}

// Config holds the API configuration values.
type Config struct {
	// Database configuration.
	DBConfig

	// API Configuration.
	APIConfig
}
