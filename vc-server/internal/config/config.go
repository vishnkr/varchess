package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const (
	dbHost = "DB_HOST"
	dbPort = "DB_PORT"
	dbUser = "DB_USER"
	dbName = "DB_NAME"
	dbPassword = "DB_PASSWORD"

	serverHost = "SERVER_HOST"
	serverPort = "SERVER_PORT"
	EnvKey = "ENVIRONMENT"
	jwtSecret = "JWT_SECRET_KEY"
	
)

type Config struct{
	ServerHost string
	ServerPort string
	DB DBConfig
	JWTSecret string
}

type DBConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

type Validate struct {
	*validator.Validate
}

func getDBConfig() (DBConfig,error){
	dbConfig := DBConfig{
		Host : os.Getenv(dbHost),
		User: os.Getenv(dbUser),
		Password: os.Getenv(dbPassword),
		Port: os.Getenv(dbPort),
		Name: os.Getenv(dbName),
	}
	validate := validator.New()
	if err := validate.Struct(dbConfig); err != nil {
		return DBConfig{},fmt.Errorf("missing database env var: %v", err)
	}
	return dbConfig, nil
}

func Load(file string) (*Config, error) {
	err := godotenv.Load(file)
	if err != nil {
		env := os.Getenv(EnvKey)
		if env==""{
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	serverPort := os.Getenv(serverPort)
	dbConfig ,err := getDBConfig()
	if err!=nil{
		return nil,err
	}
	jwtSecretValue := os.Getenv(jwtSecret)
	if jwtSecretValue == "" {
		return nil, fmt.Errorf("missing JWT_SECRET environment variable")
	}
	return &Config{
		ServerPort: serverPort,
		DB: dbConfig,
	},nil
}
