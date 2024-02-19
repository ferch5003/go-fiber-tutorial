package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/sql"
	"github.com/jmoiron/sqlx"
)

const (
	_getAllUsersStmt = `SELECT first_name, last_name, email FROM users;`
	_getUserByIDStmt = `SELECT first_name, last_name, email FROM users WHERE id = ?;`
	_saveUserStmt    = `INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?);`
	_updateUserStmt  = `UPDATE users SET %s WHERE id = ?;`
	_deleteUserStmt  = `DELETE FROM users WHERE id = ?;`
)

type Repository interface {
	// GetAll obtain all users from the database.
	GetAll(ctx context.Context) ([]domain.User, error)

	// Get obtain one User by ID.
	Get(ctx context.Context, id int) (domain.User, error)

	// Save a new User into the database.
	Save(ctx context.Context, user domain.User) (int, error)

	// Update data from the User.
	Update(ctx context.Context, user domain.User) error

	// Delete the User from the database.
	Delete(ctx context.Context, id int) error
}

type repository struct {
	conn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) Repository {
	return &repository{conn: conn}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0)

	if err := r.conn.SelectContext(ctx, &users, _getAllUsersStmt); err != nil {
		return make([]domain.User, 0), err
	}

	return users, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.User, error) {
	var user domain.User

	if err := r.conn.GetContext(ctx, &user, _getUserByIDStmt, id); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *repository) Save(ctx context.Context, user domain.User) (int, error) {
	tx, err := r.conn.Beginx()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PreparexContext(ctx, _saveUserStmt)
	if err != nil {
		return 0, err
	}

	defer func() {
		err = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return 0, rollbackErr
		}

		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (r *repository) Update(ctx context.Context, user domain.User) error {
	columns := []string{"first_name", "last_name", "email"}

	tx, err := r.conn.Beginx()
	if err != nil {
		return err
	}

	dynamicQuery, values := sql.DynamicQuery(columns, user)
	if len(dynamicQuery) <= 0 {
		return errors.New("no rows is going to be updated. User is empty")
	}

	values = append(values, user.ID)
	query := fmt.Sprintf(_updateUserStmt, dynamicQuery)

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return err
}

func (r *repository) Delete(ctx context.Context, id int) error {
	tx, err := r.conn.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, _deleteUserStmt)
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

	if err = tx.Commit(); err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return errors.New("no rows affected")
	}

	return err
}
