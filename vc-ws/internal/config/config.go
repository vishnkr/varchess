package config

import "os"

type Config struct {
    RedisAddr string
}

func LoadConfig() Config {
    return Config{
        RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
    }
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
