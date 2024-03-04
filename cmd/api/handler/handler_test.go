package handler

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/jwtauth"
	"github.com/stretchr/testify/mock"
)

var _testConfigs = &config.EnvVars{
	AppName:        "test",
	AppSecretKey:   "test",
	AppSessionType: "app",
}

var _testSessionConfigs = &jwtauth.Config{
	AppName: "test",
	Secret:  "test",
}

type errorResponse struct {
	Error string `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
}

// get token user session
func getTestUserSession() (string, error) {
	token, _, err := jwtauth.GenerateToken(1, "test", *_testSessionConfigs)
	if err != nil {
		return "", err
	}

	return token, nil
}

type sessionServiceMock struct {
	mock.Mock
}

func (ssm *sessionServiceMock) SetSession(ctx context.Context, token string, claims map[string]any) error {
	args := ssm.Called(ctx, token, claims)
	return args.Error(0)
}

func (ssm *sessionServiceMock) GetSession(ctx context.Context, token string) (map[string]string, error) {
	args := ssm.Called(ctx, token)
	return args.Get(0).(map[string]string), args.Error(1)
}
