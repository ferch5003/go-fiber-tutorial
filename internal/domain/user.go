package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id" db:"id" `
	FirstName string `json:"first_name" db:"first_name" fake:"{firstname}"`
	LastName  string `json:"last_name" db:"last_name" fake:"{lastname}"`
	Email     string `json:"email" db:"email" fake:"{email}"`
	Password  string `json:"password" db:"password"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
