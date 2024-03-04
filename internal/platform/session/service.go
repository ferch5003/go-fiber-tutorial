package session

import "context"

type Service interface {
	SetSession(ctx context.Context, token string, claims map[string]any) error
	GetSession(ctx context.Context, token string) (map[string]string, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s service) SetSession(ctx context.Context, token string, claims map[string]any) error {
	return s.repository.SetSession(ctx, token, claims)
}

func (s service) GetSession(ctx context.Context, token string) (map[string]string, error) {
	return s.repository.GetSession(ctx, token)
}
