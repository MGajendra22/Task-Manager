package user

import (
	"errors"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var err = errors.New("name and email cannot be empty")

func (u *User) Validate() error {
	if u.Name == "" || u.Email == "" {
		return err
	}

	return nil
}
