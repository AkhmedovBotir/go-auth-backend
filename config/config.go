package config

import "os"

type Config struct {
	Port         string
	JWTSecret    string
	DatabasePath string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "auth.db"
	}

	return Config{
		Port:         port,
		JWTSecret:    jwtSecret,
		DatabasePath: dbPath,
	}
}
