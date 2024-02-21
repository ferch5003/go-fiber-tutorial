package domain

type User struct {
	ID        int    `json:"id" db:"id" `
	FirstName string `json:"first_name" db:"first_name" fake:"{firstname}"`
	LastName  string `json:"last_name" db:"last_name" fake:"{lastname}"`
	Email     string `json:"email" db:"email" fake:"{email}"`
	Password  string `json:"password" db:"password"`
}
