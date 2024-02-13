package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type EnvVars struct {
	Host          string
	Port          string
	MySQLDSN      string
	MySQLUsername string
	MySQLPassword string
	MySQLDB       string
}

// dir returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func dir(envFile string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("go.mod not found")
		}

		currentDir = parent
	}

	return filepath.Join(currentDir, envFile), nil
}

func NewConfigurations() (*EnvVars, error) {
	envFilepath, err := dir(".env")
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(envFilepath); err != nil {
		return nil, err
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	mySQLDSN := os.Getenv("MYSQL_DSN")
	mySQLUsername := os.Getenv("MYSQL_USERNAME")
	mySQLPassword := os.Getenv("MYSQL_PASSWORD")
	mySQLDB := os.Getenv("MYSQL_DB")

	environment := &EnvVars{
		Host:          host,
		Port:          port,
		MySQLDSN:      mySQLDSN,
		MySQLUsername: mySQLUsername,
		MySQLPassword: mySQLPassword,
		MySQLDB:       mySQLDB,
	}

	return environment, nil
}
