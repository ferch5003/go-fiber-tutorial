package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestGetAll_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUsers := []domain.User{
		{
			FirstName: "Jhon",
			LastName:  "Smith",
			Email:     "jhon@example.com",
		},
		{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
		},
	}

	columns := []string{"first_name", "last_name", "email"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedUsers[0].FirstName, expectedUsers[0].LastName, expectedUsers[0].Email)
	rows.AddRow(expectedUsers[1].FirstName, expectedUsers[1].LastName, expectedUsers[1].Email)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	users, err := repository.GetAll(ctx)

	// Then
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, expectedUsers, users)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_FailsDueToInvalidSelect(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM users;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT first_name, last_name,
										email FROM users;\" with expected regexp \"SELECT wrong FROM users;\"`)

	var expectedUser domain.User
	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	users, err := repository.Get(ctx, 0)

	// Then
	require.Equal(t, expectedUser, users)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestGet_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 1
	expectedUser := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "jhon@example.com",
	}

	columns := []string{"first_name", "last_name", "email"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedUser.FirstName, expectedUser.LastName, expectedUser.Email)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	user, err := repository.Get(ctx, expectedUserID)

	// Then
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, expectedUser, user)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_FailsDueToInvalidSelect(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 1
	expectedUser := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "jhon@example.com",
	}

	columns := []string{"first_name", "last_name", "email"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedUser.FirstName, expectedUser.LastName, expectedUser.Email)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	user, err := repository.Get(ctx, expectedUserID)

	// Then
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, expectedUser, user)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 1
	user := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345",
	}
	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO user`)
	mock.ExpectExec(`INSERT INTO user`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.NoError(t, err)
	require.NotNil(t, userID)
	require.Equal(t, expectedUserID, userID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_FailsDueToInvalidBeginTransaction(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 0
	user := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345",
	}

	expectedError := errors.New("You have an error in your SQL syntax")

	mock.ExpectBegin().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.Equal(t, expectedUserID, userID)
	require.ErrorContains(t, err, "You have an error in your SQL syntax")
}

func TestSave_FailsDueToInvalidPreparation(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 0
	user := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345",
	}
	wrongQuery := regexp.QuoteMeta(`INSERT INTO users (first_name, last_name, email, password)) VALUES ();`)
	expectedError := errors.New(`Prepare: could not match actual sql: \"INSERT INTO users (first_name, 
										last_name, email, password) VALUES (?, ?, ?, ?);\" with expected 
										regexp \"INSERT INTO users \\(first_name, last_name, email, 
										password\\)\\) VALUES \\(\\);\"`)

	mock.ExpectBegin()
	mock.ExpectPrepare(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.Equal(t, expectedUserID, userID)
	require.ErrorContains(t, err, "Prepare: could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestSave_FailsDueToFailingExec(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 0
	user := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345",
	}

	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO user`)
	mock.ExpectExec(`INSERT INTO user`).WillReturnError(expectedError)
	mock.ExpectRollback()

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.Equal(t, expectedUserID, userID)
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
}

func TestSave_FailsDueToFailingExecWithFailingRollback(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 0
	user := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345",
	}

	expectedExecError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")
	expectedRollbackError := fmt.Errorf("insert failed: %v, unable to back: %v",
		expectedExecError, "Rollack error")

	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO user`)
	mock.ExpectExec(`INSERT INTO user`).WillReturnError(expectedExecError)
	mock.ExpectRollback().WillReturnError(expectedRollbackError)

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.Equal(t, expectedUserID, userID)
	require.ErrorContains(t, err, "insert failed")
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
	require.ErrorContains(t, err, "unable to back")
	require.ErrorContains(t, err, "Rollack error")
}

func TestSave_FailsDueToFailingCommit(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 0
	user := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345",
	}
	expectedError := errors.New("sql: transaction has already been committed or rolled back")

	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO user`)
	mock.ExpectExec(`INSERT INTO user`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.Equal(t, expectedUserID, userID)
	require.ErrorContains(t, err, "sql")
	require.ErrorContains(t, err, "transaction has already been committed or rolled back")
}
