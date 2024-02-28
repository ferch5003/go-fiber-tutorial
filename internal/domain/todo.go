package domain

type Todo struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" fake:"{word}"`
	Description string `json:"description" db:"description" fake:"{loremipsumsentence:10}"`
	Completed   bool   `json:"completed" db:"completed" fake:"{bool}"`
	UserID      int    `json:"user_id"`
}
