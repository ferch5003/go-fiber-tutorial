package seeds

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
)

// todosSeed seeds todos data.
func (s Seed) todosSeed() {
	for range 10 {
		var todo domain.Todo
		if err := gofakeit.Struct(&todo); err != nil {
			panic(err)
		}

		todo.UserID = 1

		_, err := s.todoRepository.Save(context.Background(), todo)
		if err != nil {
			panic(fmt.Sprintf("error seeding todos: %v", err))
		}
	}
}
