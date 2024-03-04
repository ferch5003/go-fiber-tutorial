package config

import (
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/files"
	"github.com/joho/godotenv"
	"os"
)

type EnvVars struct {
	// App Data.
	AppName        string
	AppSecretKey   string
	AppSessionType string // Fiber JWT or Manual JWT ("fiber" or "app").
	Host           string
	Port           string

	// MySQL Data.
	MySQLDSN      string
	MySQLUsername string
	MySQLPassword string
	MySQLDB       string

	// Redis Data.
	RedisConnection string
	RedisUsername   string
	RedisPassword   string
	RedisDB         string
}

func NewConfigurations() (*EnvVars, error) {
	envFilepath, err := files.GetFile(".env")
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(envFilepath); err != nil {
		return nil, err
	}

	appName := os.Getenv("APP_NAME")
	appSecretKey := os.Getenv("APP_SECRET_KEY")
	appSessionType := os.Getenv("APP_SESSION_TYPE")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	mySQLDSN := os.Getenv("MYSQL_DSN")
	mySQLUsername := os.Getenv("MYSQL_USERNAME")
	mySQLPassword := os.Getenv("MYSQL_PASSWORD")
	mySQLDB := os.Getenv("MYSQL_DB")

	redisConnection := os.Getenv("REDIS_CONNECTION")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisLDB := os.Getenv("REDIS")

	environment := &EnvVars{
		AppName:        appName,
		AppSecretKey:   appSecretKey,
		AppSessionType: appSessionType,
		Host:           host,
		Port:           port,

		MySQLDSN:      mySQLDSN,
		MySQLUsername: mySQLUsername,
		MySQLPassword: mySQLPassword,
		MySQLDB:       mySQLDB,

		RedisConnection: redisConnection,
		RedisUsername:   redisUsername,
		RedisPassword:   redisPassword,
		RedisDB:         redisLDB,
	}

	return environment, nil
}
