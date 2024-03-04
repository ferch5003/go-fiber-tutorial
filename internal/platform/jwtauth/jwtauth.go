package jwtauth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Config struct {
	AppName string
	Secret  string
}

func GenerateToken(id int, name string, cfg Config) (string, map[string]any, error) {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return "", nil, err
	}

	claims := jwt.MapClaims{
		"iss":  cfg.AppName,
		"sub":  id,
		"name": name,
		"exp":  time.Now().In(location).Add(72 * time.Hour).Unix(),
		"iat":  time.Now().In(location).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", nil, err
	}

	return t, claims, nil
}
