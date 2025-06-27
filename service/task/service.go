package task

import (
	"Task_Manager/model/task"
	storePkg "Task_Manager/store/task"
)

type Service struct {
	str *storePkg.Store
}

func NewService(s *storePkg.Store) *Service {
	return &Service{str: s}
}

func (s *Service) Create(t task.Task) (task.Task, error) {
	if err := t.Validate(); err != nil {
		return t, err
	}
	return s.str.Create(t)
}

func (s *Service) GetTask(id int) (task.Task, error) {
	return s.str.GetByID(id)
}

func (s *Service) Complete(id int) error {
	return s.str.Complete(id)
}

func (s *Service) Delete(id int) error {
	return s.str.Delete(id)
}

func (s *Service) All() ([]task.Task, error) {
	return s.str.GetAll()
}

func (s *Service) GetTasksByUserID(userid int) ([]task.Task, error) {
	return s.str.GetTasksByUserID(userid)
}
