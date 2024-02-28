package todo

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
)

type Service interface {
	// GetAll obtain all todos from the database of specific user.
	GetAll(ctx context.Context, userID int) ([]domain.Todo, error)

	// Get obtain one Todo by ID.
	Get(ctx context.Context, id int) (domain.Todo, error)

	// Save a new Todo into the database.
	Save(ctx context.Context, todo domain.Todo) (domain.Todo, error)

	// Completed change the completed state to true.
	Completed(ctx context.Context, id int) error

	// Delete the Todo from the database.
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s service) GetAll(ctx context.Context, userID int) ([]domain.Todo, error) {
	return s.repository.GetAll(ctx, userID)
}

func (s service) Get(ctx context.Context, id int) (domain.Todo, error) {
	return s.repository.Get(ctx, id)
}

func (s service) Save(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	id, err := s.repository.Save(ctx, todo)
	if err != nil {
		return domain.Todo{}, err
	}

	todo.ID = id

	return todo, nil
}

func (s service) Completed(ctx context.Context, id int) error {
	return s.repository.Completed(ctx, id)
}

func (s service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
