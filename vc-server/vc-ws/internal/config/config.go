package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	serverHost = "SERVER_HOST"
	serverPort = "SERVER_PORT"
	dbUriEnv = "DB_URI"
	dbName = "DB_NAME"
	EnvKey = "ENVIRONMENT"
	jwtSecret = "JWT_SECRET_KEY"
	redisAddr = "REDIS_ADDR"
	redisPassword = "REDIS_PASSWORD"
	
)

type Config struct {
    RedisAddr string
    ServerHost string
    ServerPort string
}

func Load(file string) (*Config,error) {
    err := godotenv.Load(file)
	if err != nil {
		env := os.Getenv(EnvKey)
		if env==""{
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

    return &Config{
        RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
        ServerHost: os.Getenv("SERVER_HOST"),
        ServerPort: os.Getenv("SERVER_PORT"),
    },nil
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
