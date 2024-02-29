package handler

import "github.com/ferch5003/go-fiber-tutorial/internal/platform/jwtauth"

var _testSessionConfigs = &jwtauth.Config{
	AppName: "test",
	Secret:  "test",
}

type errorResponse struct {
	Error string `json:"error"`
}

// get token user session
func getTestUserSession() (string, error) {
	token, err := jwtauth.GenerateToken(1, "test", *_testSessionConfigs)
	if err != nil {
		return "", err
	}

	return token, nil
}
