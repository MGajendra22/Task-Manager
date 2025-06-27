package user

import (
	"Task_Manager/model/user"
	storePkg "Task_Manager/store/user"
)

type UserService struct {
	store *storePkg.UserStore
}

func NewUserService(store *storePkg.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Create(u user.User) (user.User, error) {
	if err := u.Validate(); err != nil {
		return u, err
	}
	return s.store.Create(u)
}

func (s *UserService) Get(id int) (user.User, error) {
	return s.store.GetByID(id)
}

func (s *UserService) Delete(id int) error {
	return s.store.Delete(id)
}

func (s *UserService) All() ([]user.User, error) {
	return s.store.GetAll()
}
