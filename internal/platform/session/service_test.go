package session

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) SetSession(ctx context.Context, token string, claims map[string]any) error {
	args := m.Called(ctx, token, claims)
	return args.Error(0)
}

func (m *mockRepository) GetSession(ctx context.Context, token string) (map[string]string, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestServiceSetSession_Successful(t *testing.T) {
	// Given
	expectedToken := "token_db"
	claims := map[string]any{
		"iss":  "test",
		"sub":  1,
		"name": "John Doe",
		"exp":  time.Now().Add(72 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	mr := new(mockRepository)
	mr.On("SetSession", mock.Anything, expectedToken, claims).Return(nil)

	service := NewService(mr)

	// When
	err := service.SetSession(context.Background(), expectedToken, claims)

	// Then
	require.NoError(t, err)
}

func TestServiceGetSession_Successful(t *testing.T) {
	// Given
	expectedToken := "token_db"
	expectedClaims := map[string]string{
		"iss":  "test",
		"sub":  "1",
		"name": "John Doe",
		"exp":  fmt.Sprintf("%v", time.Now().Add(72*time.Hour).Unix()),
		"iat":  fmt.Sprintf("%v", time.Now().Unix()),
	}

	mr := new(mockRepository)
	mr.On("GetSession", mock.Anything, expectedToken).Return(expectedClaims, nil)

	service := NewService(mr)

	// When
	claims, err := service.GetSession(context.Background(), expectedToken)

	// Then
	require.NoError(t, err)
	require.EqualValues(t, expectedClaims, claims)
}
