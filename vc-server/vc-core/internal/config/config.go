package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const (
	serverHost      = "SERVER_HOST"
	serverPort      = "SERVER_PORT"
	dbUriEnv        = "DB_URI"
	dbName          = "DB_NAME"
	EnvKey          = "ENVIRONMENT"
	jwtSecret       = "JWT_SECRET_KEY"
	redisAddr       = "REDIS_ADDR"
	redisPassword   = "REDIS_PASSWORD"
	clientURLEnv    = "CLIENT_URL"
	publicAPIURLEnv = "PUBLIC_API_URL"
	githubClientID  = "GITHUB_CLIENT_ID"
	githubSecret    = "GITHUB_CLIENT_SECRET"
	googleClientID  = "GOOGLE_CLIENT_ID"
	googleSecret    = "GOOGLE_CLIENT_SECRET"
)

type Config struct {
	ServerHost         string
	ServerPort         string
	DBURI              string
	JWTSecret          string
	DBName             string
	RedisAddr          string
	RedisPassword      string
	ClientURL          string
	PublicAPIURL       string
	GitHubClientID     string
	GitHubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
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

	githubID := firstEnv(githubClientID, "GITHUB_OAUTH_CLIENT_ID")
	githubSec := firstEnv(githubSecret, "GITHUB_OAUTH_CLIENT_SECRET")
	googleID := firstEnv(googleClientID, "GOOGLE_OAUTH_CLIENT_ID")
	googleSec := firstEnv(googleSecret, "GOOGLE_OAUTH_CLIENT_SECRET")

	return &Config{
		ServerHost:         os.Getenv(serverHost),
		ServerPort:         serverPort,
		DBURI:              dbUri,
		DBName:             dbName,
		RedisAddr:          redisAddr,
		RedisPassword:      redisPass,
		JWTSecret:          jwtSecretValue,
		ClientURL:          firstEnv(clientURLEnv, "FRONTEND_URL"),
		PublicAPIURL:       firstEnv(publicAPIURLEnv, "API_URL"),
		GitHubClientID:     githubID,
		GitHubClientSecret: githubSec,
		GoogleClientID:     googleID,
		GoogleClientSecret: googleSec,
	}, nil
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
