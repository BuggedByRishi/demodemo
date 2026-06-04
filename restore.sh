cat > internal/config/config.go << 'EOF'
package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	Env         string
}

func Load() Config {
	host     := getenvOr("DB_HOST", "localhost")
	port     := getenvOr("DB_PORT", "5432")
	user     := mustGetenv("DB_USER")
	password := mustGetenv("DB_PASSWORD")
	dbname   := mustGetenv("DB_NAME")
	sslmode  := getenvOr("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	return Config{
		DatabaseURL: dsn,
		Port:        getenvOr("PORT", "8080"),
		Env:         getenvOr("APP_ENV", "development"),
	}
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return v
}

func getenvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
EOF