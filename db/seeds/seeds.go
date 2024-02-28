package seeds

import (
	"github.com/danvergara/seeder"
	"github.com/ferch5003/go-fiber-tutorial/internal/todo"
	"github.com/ferch5003/go-fiber-tutorial/internal/user"
)

// Seed struct.
type Seed struct {
	userRepository user.Repository
	todoRepository todo.Repository
}

// NewSeed return a Seed with a pool of connection to a dabase.
func NewSeed(userRepository user.Repository, todoRepository todo.Repository) Seed {
	return Seed{
		userRepository: userRepository,
		todoRepository: todoRepository,
	}
}

func Execute(s Seed) error {
	return seeder.Execute(s)
}

func (s Seed) PopulateDB() {
	s.usersSeed()
	s.todosSeed()
}
