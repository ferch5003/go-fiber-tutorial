package seeds

import (
	"github.com/danvergara/seeder"
	"github.com/ferch5003/go-fiber-tutorial/internal/user"
)

// Seed struct.
type Seed struct {
	userRepository user.Repository
}

// NewSeed return a Seed with a pool of connection to a dabase.
func NewSeed(userRepository user.Repository) Seed {
	return Seed{
		userRepository: userRepository,
	}
}

func Execute(s Seed) error {
	return seeder.Execute(s)
}
