package user

import "Task_Manager/model/user"

type UserServiceInterface interface {
	Create(u user.User) (user.User, error)
	Get(id int) (user.User, error)
	Delete(id int) error
	All() ([]user.User, error)
}
