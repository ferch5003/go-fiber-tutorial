package user

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
)

type Service interface {
	// GetAll obtain all users.
	GetAll(ctx context.Context) ([]domain.User, error)

	// Get obtain one User by ID.
	Get(ctx context.Context, id int) (domain.User, error)

	// Save a new User.
	Save(ctx context.Context, user domain.User) (domain.User, error)

	// Update data from the User.
	Update(ctx context.Context, user domain.User) (domain.User, error)

	// Delete the User.
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

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.repository.GetAll(ctx)
}

func (s service) Get(ctx context.Context, id int) (domain.User, error) {
	return s.repository.Get(ctx, id)
}

func (s service) Save(ctx context.Context, user domain.User) (domain.User, error) {
	id, err := s.repository.Save(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	user.ID = id

	return user, nil
}

func (s service) Update(ctx context.Context, user domain.User) (domain.User, error) {
	if err := s.repository.Update(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
