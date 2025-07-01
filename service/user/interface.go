package user

import "Task_Manager/model/user"

type UserStoreInterface interface {
	CreateUser(u user.User) (user.User, error)
	GetByIDUser(id int) (user.User, error)
	DeleteUser(id int) error
	GetAllUser() ([]user.User, error)
}
