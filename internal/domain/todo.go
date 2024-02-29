package domain

type Todo struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" fake:"{word}" validate:"required"`
	Description string `json:"description" db:"description" fake:"{loremipsumsentence:10}" validate:"required"`
	Completed   bool   `json:"completed" db:"completed" fake:"{bool}"`
	UserID      int    `json:"user_id" db:"user_id"`
}
