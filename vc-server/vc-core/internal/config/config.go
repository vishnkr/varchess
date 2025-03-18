package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
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

type Config struct{
	ServerHost string
	ServerPort string
	DBURI string
	JWTSecret string
	DBName string
	RedisAddr string
	RedisPassword string
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

/*func getDBConfig() (DBConfig,error){
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
}*/

func Load(file string) (*Config, error) {
	err := godotenv.Load(file)
	if err != nil {
		env := os.Getenv(EnvKey)
		if env==""{
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	serverPort := os.Getenv(serverPort)
	dbUri := os.Getenv(dbUriEnv)
	dbName := os.Getenv(dbName)
	redisAddr := os.Getenv(redisAddr)
	redisPass := os.Getenv(redisPassword) 
	/*dbConfig ,err := getDBConfig()
	if err!=nil{
		return nil,err
	}*/
	jwtSecretValue := os.Getenv(jwtSecret)
	if jwtSecretValue == "" {
		return nil, fmt.Errorf("missing JWT_SECRET_KEY environment variable")
	}
	return &Config{
		ServerPort: serverPort,
		DBURI: dbUri,
		DBName: dbName,
		RedisAddr: redisAddr,
		RedisPassword: redisPass,
	},nil
}
