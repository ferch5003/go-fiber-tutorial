package todo

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

	expectedUserID := 1
	expectedTodos := []domain.Todo{
		{
			ID:          1,
			Title:       "Lorem",
			Description: "Ipsum",
			Completed:   false,
		},
		{
			ID:          2,
			Title:       "Lorem Ipsum",
			Description: "FLCL",
			Completed:   true,
		},
	}

	columns := []string{"id", "title", "description", "completed"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedTodos[0].ID, expectedTodos[0].Title, expectedTodos[0].Description, expectedTodos[0].Completed)
	rows.AddRow(expectedTodos[1].ID, expectedTodos[1].Title, expectedTodos[1].Description, expectedTodos[1].Completed)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	todos, err := repository.GetAll(ctx, expectedUserID)

	// Then
	require.NoError(t, err)
	require.NotNil(t, todos)
	require.Equal(t, expectedTodos, todos)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetAll_FailsDueToInvalidSelect(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	expectedUserID := 1

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM todos;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT id, title, description, completed 
										FROM todos
										INNER JOIN users ON
										users.id = todos.user_id
										WHERE user_id = ?;\" with expected regexp 
										\"SELECT wrong FROM todos;\"`)

	expectedTodos := make([]domain.Todo, 0)
	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	todos, err := repository.GetAll(ctx, expectedUserID)

	// Then
	require.Equal(t, expectedTodos, todos)
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

	expectedTodo := domain.Todo{
		ID:          1,
		Title:       "Lorem",
		Description: "Ipsum",
		Completed:   false,
	}

	columns := []string{"id", "title", "description", "completed"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(expectedTodo.ID, expectedTodo.Title, expectedTodo.Description, expectedTodo.Completed)
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	repository := NewRepository(dbx)

	// When
	todo, err := repository.Get(ctx, expectedTodo.ID)

	// Then
	require.NoError(t, err)
	require.NotNil(t, todo)
	require.Equal(t, expectedTodo, todo)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGet_FailsDueToInvalidGet(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	wrongQuery := regexp.QuoteMeta("SELECT wrong FROM todos;")
	expectedError := errors.New(`Query: could not match actual sql: \"SELECT id, title, description, completed 
										FROM todos
										WHERE id = ?\" with expected regexp 
										\"SELECT wrong FROM todos;\"`)
	expectedTodo := domain.Todo{}

	mock.ExpectQuery(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	todo, err := repository.Get(ctx, 0)

	// Then
	require.Equal(t, expectedTodo, todo)
	require.ErrorContains(t, err, "Query")
	require.ErrorContains(t, err, "could not match actual sql")
	require.ErrorContains(t, err, "with expected regexp")
}

func TestRepositoryDelete_Successful(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	deletedTodoID := 1
	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM todos`)
	mock.ExpectExec(`DELETE FROM todos`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedTodoID)

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

	deletedTodoID := 0
	expectedError := errors.New("You have an error in your SQL syntax")

	mock.ExpectBegin().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedTodoID)

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

	deletedTodoID := 0
	wrongQuery := regexp.QuoteMeta(`DELETE FROM todos WHERE id = ();`)
	expectedError := errors.New(`Prepare: could not match actual sql: \"DELETE FROM todos WHERE id = ?;\" 
										with expected regexp \"DELETE FROM todos WHERE id = ();\"`)

	mock.ExpectBegin()
	mock.ExpectPrepare(wrongQuery).WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedTodoID)

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

	deletedTodoID := 0
	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM todos`)
	mock.ExpectExec(`DELETE FROM todos`).WillReturnError(expectedError)
	mock.ExpectRollback()

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedTodoID)

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

	deletedTodoID := 0
	expectedExecError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")
	expectedRollbackError := fmt.Errorf("delete failed: %v, unable to back: %v",
		expectedExecError, "Rollack error")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM todos`)
	mock.ExpectExec(`DELETE FROM todos`).WillReturnError(expectedExecError)
	mock.ExpectRollback().WillReturnError(expectedRollbackError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedTodoID)

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

	deletedTodoID := 0
	expectedError := errors.New("sql: transaction has already been committed or rolled back")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM todos`)
	mock.ExpectExec(`DELETE FROM todos`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(expectedError)

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedTodoID)

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

	deletedTodoID := 0
	expectedError := errors.New("no rows affected")

	mock.ExpectBegin()
	mock.ExpectPrepare(`DELETE FROM todos`)
	mock.ExpectExec(`DELETE FROM todos`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repository := NewRepository(dbx)

	// When
	err = repository.Delete(ctx, deletedTodoID)

	// Then
	require.ErrorContains(t, err, expectedError.Error())
}
