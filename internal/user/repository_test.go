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

func TestRepositoryGetAll_Successful(t *testing.T) {
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

func TestRepositoryGetAll_FailsDueToInvalidSelect(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM users;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT first_name, last_name,
										email FROM users;\" with expected regexp \"SELECT wrong FROM users;\"`)

	expectedUsers := make([]domain.User, 0)
	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	users, err := repository.GetAll(ctx)

	// Then
	require.Equal(t, expectedUsers, users)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositoryGet_Successful(t *testing.T) {
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

func TestRepositoryGet_FailsDueToInvalidGet(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM users;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT first_name, last_name,
										email FROM users WHERE id = ?;\" with expected regexp \"SELECT 
										wrong FROM users;\"`)

	expectedUser := domain.User{}

	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	user, err := repository.Get(ctx, 0)

	// Then
	require.Equal(t, expectedUser, user)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositoryGetByEmail_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUser := domain.User{
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "jhon@example.com",
		Password:  "12345678",
	}

	columns := []string{"id", "first_name", "last_name", "email", "password"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Password)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	user, err := repository.GetByEmail(ctx, expectedUser.Email)

	// Then
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, expectedUser, user)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetByEmail_FailsDueToInvalidGet(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM users;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT first_name, last_name, 
										email, password FROM users WHERE email = ?;\" with expected regexp
										\"SELECT wrong FROM users;\"`)

	expectedUser := domain.User{}

	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	user, err := repository.Get(ctx, 0)

	// Then
	require.Equal(t, expectedUser, user)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositorySave_Successful(t *testing.T) {
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
	mock.ExpectPrepare(`INSERT INTO users`)
	mock.ExpectExec(`INSERT INTO users`).WillReturnResult(sqlmock.NewResult(1, 1))
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

func TestRepositorySave_FailsDueToInvalidBeginTransaction(t *testing.T) {
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

func TestRepositorySave_FailsDueToInvalidPreparation(t *testing.T) {
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

func TestRepositorySave_FailsDueToFailingExec(t *testing.T) {
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
	mock.ExpectPrepare(`INSERT INTO users`)
	mock.ExpectExec(`INSERT INTO users`).WillReturnError(expectedError)
	mock.ExpectRollback()

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.Equal(t, expectedUserID, userID)
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
}

func TestRepositorySave_FailsDueToFailingExecWithFailingRollback(t *testing.T) {
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
	mock.ExpectPrepare(`INSERT INTO users`)
	mock.ExpectExec(`INSERT INTO users`).WillReturnError(expectedExecError)
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

func TestRepositorySave_FailsDueToFailingCommit(t *testing.T) {
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
	mock.ExpectPrepare(`INSERT INTO users`)
	mock.ExpectExec(`INSERT INTO users`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	userID, err := repository.Save(ctx, user)

	// Then
	require.Equal(t, expectedUserID, userID)
	require.ErrorContains(t, err, "sql")
	require.ErrorContains(t, err, "transaction has already been committed or rolled back")
}

func TestRepositoryUpdate_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	user := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
	}
	mock.ExpectBegin()
	mock.ExpectPrepare(`UPDATE users`)
	mock.ExpectExec(`UPDATE users`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := NewRepository(dbx)

	// When
	err = repository.Update(ctx, user)

	// Then
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryUpdate_FailsDueToInvalidBeginTransaction(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	user := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	expectedError := errors.New("You have an error in your SQL syntax")

	mock.ExpectBegin().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Update(ctx, user)

	// Then
	require.ErrorContains(t, err, "You have an error in your SQL syntax")
}

func TestRepositoryUpdate_FailsDueToNoneColumnsToUpdate(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	user := domain.User{
		ID: 1,
	}

	mock.ExpectBegin()

	repository := NewRepository(dbx)

	// When
	err = repository.Update(ctx, user)

	// Then
	require.ErrorContains(t, err, "no rows is going to be updated")
	require.ErrorContains(t, err, "User is empty")
}

func TestRepositoryUpdate_FailsDueToInvalidPreparation(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	user := domain.User{
		ID:        1,
		FirstName: "Jhon",
	}
	wrongQuery := regexp.QuoteMeta(`UPDATE users SET  first_name = ?? WHERE id = ?;`)
	expectedError := errors.New(`Prepare: could not match actual sql: \"UPDATE users SET first_name = ? 
										WHERE id = ?;\" with expected  regexp \"UPDATE users SET  first_name = ?? 
										WHERE id = ?;"`)

	mock.ExpectBegin()
	mock.ExpectPrepare(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Update(ctx, user)

	// Then
	require.ErrorContains(t, err, "Prepare: could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositoryUpdate_FailsDueToFailingExec(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	user := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	mock.ExpectBegin()
	mock.ExpectPrepare(`UPDATE users`)
	mock.ExpectExec(`UPDATE users`).WillReturnError(expectedError)
	mock.ExpectRollback()

	repository := NewRepository(dbx)

	// When
	err = repository.Update(ctx, user)

	// Then
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
}

func TestRepositoryUpdate_FailsDueToFailingExecWithFailingRollback(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	user := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	expectedExecError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")
	expectedRollbackError := fmt.Errorf("insert failed: %v, unable to back: %v",
		expectedExecError, "Rollack error")

	mock.ExpectBegin()
	mock.ExpectPrepare(`UPDATE users`)
	mock.ExpectExec(`UPDATE users`).WillReturnError(expectedExecError)
	mock.ExpectRollback().WillReturnError(expectedRollbackError)

	repository := NewRepository(dbx)

	// When
	err = repository.Update(ctx, user)

	// Then
	require.ErrorContains(t, err, "insert failed")
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
	require.ErrorContains(t, err, "unable to back")
	require.ErrorContains(t, err, "Rollack error")
}

func TestRepositoryUpdate_FailsDueToFailingCommit(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	user := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
	}
	expectedError := errors.New("sql: transaction has already been committed or rolled back")

	mock.ExpectBegin()
	mock.ExpectPrepare(`UPDATE users`)
	mock.ExpectExec(`UPDATE users`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Update(ctx, user)

	// Then
	require.ErrorContains(t, err, "sql")
	require.ErrorContains(t, err, "transaction has already been committed or rolled back")
}

func TestRepositoryDelete_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedUserID := 1
	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM users`)
	mock.ExpectExec(`DELETE FROM users`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedUserID)

	// Then
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryDelete_FailsDueToInvalidBeginTransaction(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedUserID := 0
	expectedError := errors.New("You have an error in your SQL syntax")

	mock.ExpectBegin().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedUserID)

	// Then
	require.ErrorContains(t, err, "You have an error in your SQL syntax")
}

func TestRepositoryDelete_FailsDueToInvalidPreparation(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedUserID := 0
	wrongQuery := regexp.QuoteMeta(`DELETE FROM users WHERE id = ();`)
	expectedError := errors.New(`Prepare: could not match actual sql: \"DELETE FROM users WHERE id = ?;\" 
										with expected regexp \"DELETE FROM users WHERE id = ();\"`)

	mock.ExpectBegin()
	mock.ExpectPrepare(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedUserID)

	// Then
	require.ErrorContains(t, err, "Prepare: could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositoryDelete_FailsDueToFailingExec(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedUserID := 0
	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM users`)
	mock.ExpectExec(`DELETE FROM users`).WillReturnError(expectedError)
	mock.ExpectRollback()

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedUserID)

	// Then
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
}

func TestRepositoryDelete_FailsDueToFailingExecWithFailingRollback(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedUserID := 0
	expectedExecError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")
	expectedRollbackError := fmt.Errorf("delete failed: %v, unable to back: %v",
		expectedExecError, "Rollack error")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM users`)
	mock.ExpectExec(`DELETE FROM users`).WillReturnError(expectedExecError)
	mock.ExpectRollback().WillReturnError(expectedRollbackError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedUserID)

	// Then
	require.ErrorContains(t, err, "delete failed")
	require.ErrorContains(t, err, "Error Code: 1136")
	require.ErrorContains(t, err, "Column count doesn't match value count at row 1")
	require.ErrorContains(t, err, "unable to back")
	require.ErrorContains(t, err, "Rollack error")
}

func TestRepositoryDelete_FailsDueToFailingCommit(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedUserID := 0
	expectedError := errors.New("sql: transaction has already been committed or rolled back")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM users`)
	mock.ExpectExec(`DELETE FROM users`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedUserID)

	// Then
	require.ErrorContains(t, err, "sql")
	require.ErrorContains(t, err, "transaction has already been committed or rolled back")
}

func TestRepositoryDelete_FailsDueToNoRowsAffected(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedUserID := 0
	expectedError := errors.New("no rows affected")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM users`)
	mock.ExpectExec(`DELETE FROM users`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedUserID)

	// Then
	require.ErrorContains(t, err, expectedError.Error())
}
