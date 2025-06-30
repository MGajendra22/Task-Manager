package user

import (
	"Task_Manager/model/user"
)

type UserStoreInterface interface {
	CreateUser(u user.User) (user.User, error)
	GetByIDUser(id int) (user.User, error)
	DeleteUser(id int) error
	GetAllUser() ([]user.User, error)
}

type UserService struct {
	store UserStoreInterface
}

func NewUserService(store UserStoreInterface) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Create(u user.User) (user.User, error) {
	if err := u.Validate(); err != nil {
		return u, err
	}

	return s.store.CreateUser(u)
}

func (s *UserService) Get(id int) (user.User, error) {
	return s.store.GetByIDUser(id)
}

func (s *UserService) Delete(id int) error {
	return s.store.DeleteUser(id)
}

func (s *UserService) All() ([]user.User, error) {
	return s.store.GetAllUser()
}
