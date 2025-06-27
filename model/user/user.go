package user

import "fmt"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) Validate() error {
	if u.Name == "" || u.Email == "" {
		return fmt.Errorf("name and email cannot be empty")
	}
	return nil
}
