package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRepositorySetSession_Successful(t *testing.T) {
	// Give
	db, mock := redismock.NewClientMock()

	token := "token_db"
	key := fmt.Sprintf("user:%s", token)
	expectedVal := map[string]any{
		"iss":  "test",
		"sub":  1,
		"name": "John Doe",
		"exp":  time.Now().Add(72 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	mock.ExpectExists(key).SetVal(0) // 0 indicates session not exists
	mock.ExpectHSet(key, expectedVal).SetVal(1)

	repository := NewRepository(db)

	// When
	err := repository.SetSession(context.Background(), token, expectedVal)

	// Then
	require.NoError(t, err)
}

func TestRepositorySetSession_FailsDueToNotExpectedMap(t *testing.T) {
	// Give
	db, mock := redismock.NewClientMock()

	token := "token_db"
	key := fmt.Sprintf("user:%s", token)
	expectedVal := map[string]any{
		"not_valid": "value",
	}
	expectedError := errors.New("not expected value")

	mock.ExpectExists(key).SetVal(0) // 0 indicates session not exists
	mock.ExpectHSet(key, expectedVal).SetErr(expectedError)

	repository := NewRepository(db)

	// When
	err := repository.SetSession(context.Background(), token, expectedVal)

	// Then
	require.Equal(t, expectedError, err)
}

func TestRepositoryGetSession_Successful(t *testing.T) {
	// Give
	db, mock := redismock.NewClientMock()

	token := "token_db"
	key := fmt.Sprintf("user:%s", token)
	expectedVal := map[string]string{
		"iss":  "test",
		"sub":  "1",
		"name": "John Doe",
		"exp":  fmt.Sprintf("%v", time.Now().Add(72*time.Hour).Unix()),
		"iat":  fmt.Sprintf("%v", time.Now().Unix()),
	}
	mock.ExpectHGetAll(key).SetVal(expectedVal)

	repository := NewRepository(db)

	// When
	session, err := repository.GetSession(context.Background(), token)

	// Then
	require.NoError(t, err)
	require.EqualValues(t, expectedVal, session)
}

func TestRepositoryGetSession_FailsDueToInvalidKey(t *testing.T) {
	// Give
	db, mock := redismock.NewClientMock()

	token := "token_db"
	key := fmt.Sprintf("user:%s", token)
	expectedErr := errors.New("args not `DeepEqual`, expectation: 'user:token_db', but gave: 'user:other_token'")
	mock.ExpectHGetAll(key).SetErr(expectedErr)

	repository := NewRepository(db)

	// When
	_, err := repository.GetSession(context.Background(), "other_token")

	// Then
	require.Equal(t, expectedErr, err)
}
