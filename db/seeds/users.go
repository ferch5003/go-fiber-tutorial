package seeds

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
)

// UsersSeed seeds roles data.
func (s Seed) UsersSeed() {
	for range 5 {
		var user domain.User
		if err := gofakeit.Struct(&user); err != nil {
			panic(err)
		}

		user.Password = "12345678"

		_, err := s.userRepository.Save(context.Background(), user)
		if err != nil {
			panic(fmt.Sprintf("error seeding users: %v", err))
		}
	}
}
