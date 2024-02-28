package seeds

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// usersSeed seeds user data.
func (s Seed) usersSeed() {
	for range 5 {
		var user domain.User
		if err := gofakeit.Struct(&user); err != nil {
			panic(err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
		if err != nil {
			panic(fmt.Sprintf("error seeding users: %v", err))
		}

		user.Password = string(hashedPassword)

		_, err = s.userRepository.Save(context.Background(), user)
		if err != nil {
			panic(fmt.Sprintf("error seeding users: %v", err))
		}
	}
}
