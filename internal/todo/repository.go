package todo

import (
	"context"
	"errors"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	_getAllTodosStmt = `SELECT id, title, description, completed 
						FROM todos
						INNER JOIN users ON
						users.id = todos.user_id
						WHERE user_id = ?;`
	_getTodoStmt = `SELECT id, title, description, completed 
					FROM todos
					WHERE id = ?;`
	_saveTodoStmt = `INSERT INTO todos (title, description, user_id) 
								VALUES (?, ?, ?);`
	_updateTodoCompletedStmt = `UPDATE todos 
								SET completed = true
								WHERE id = ?;`
	_deleteTodoStmt = `DELETE FROM todos WHERE id = ?;`
)

type Repository interface {
	// GetAll obtain all todos from the database of specific user.
	GetAll(ctx context.Context, userID int) ([]domain.Todo, error)

	// Get obtain one Todo by ID.
	Get(ctx context.Context, id int) (domain.Todo, error)

	// Save a new Todo into the database.
	Save(ctx context.Context, todo domain.Todo) (int, error)

	// Completed change the completed state to true.
	Completed(ctx context.Context, id int) error

	// Delete the Todo from the database.
	Delete(ctx context.Context, id int) error
}

type repository struct {
	conn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) Repository {
	return &repository{
		conn: conn,
	}
}

func (r repository) GetAll(ctx context.Context, userID int) ([]domain.Todo, error) {
	todos := make([]domain.Todo, 0)

	if err := r.conn.SelectContext(ctx, &todos, _getAllTodosStmt, userID); err != nil {
		return make([]domain.Todo, 0), err
	}

	return todos, nil
}

func (r repository) Get(ctx context.Context, id int) (domain.Todo, error) {
	var todo domain.Todo

	if err := r.conn.GetContext(ctx, &todo, _getTodoStmt, id); err != nil {
		return domain.Todo{}, err
	}

	return todo, nil
}

func (r repository) Save(ctx context.Context, todo domain.Todo) (int, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PreparexContext(ctx, _saveTodoStmt)
	if err != nil {
		return 0, err
	}

	defer func() {
		err = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, todo.Title, todo.Description, todo.UserID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return 0, rollbackErr
		}

		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (r repository) Completed(ctx context.Context, id int) error {
	tx, err := r.conn.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, _updateTodoCompletedStmt)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return err
}

func (r repository) Delete(ctx context.Context, id int) error {
	tx, err := r.conn.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, _deleteTodoStmt)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return errors.New("no rows affected")
	}

	return nil
}
