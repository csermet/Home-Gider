package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "gider"),
		DBPassword: getEnv("DB_PASSWORD", "gider123"),
		DBName:     getEnv("DB_NAME", "gider_db"),
		JWTSecret:  getEnv("JWT_SECRET", "replace-with-secure-secret"),
		Port:       getEnv("PORT", "3000"),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user='" + c.DBUser + "'" +
		" password='" + c.DBPassword + "'" +
		" dbname='" + c.DBName + "'" +
		" sslmode=disable"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
