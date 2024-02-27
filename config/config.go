package config

import (
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/files"
	"github.com/joho/godotenv"
	"os"
)

type EnvVars struct {
	AppName       string
	AppSecretKey  string
	Host          string
	Port          string
	MySQLDSN      string
	MySQLUsername string
	MySQLPassword string
	MySQLDB       string
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
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	mySQLDSN := os.Getenv("MYSQL_DSN")
	mySQLUsername := os.Getenv("MYSQL_USERNAME")
	mySQLPassword := os.Getenv("MYSQL_PASSWORD")
	mySQLDB := os.Getenv("MYSQL_DB")

	environment := &EnvVars{
		AppName:       appName,
		AppSecretKey:  appSecretKey,
		Host:          host,
		Port:          port,
		MySQLDSN:      mySQLDSN,
		MySQLUsername: mySQLUsername,
		MySQLPassword: mySQLPassword,
		MySQLDB:       mySQLDB,
	}

	return environment, nil
}
