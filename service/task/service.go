package task

import (
	"Task_Manager/model/task"
	"fmt"
)

type TaskService struct {
	str            TaskStoreInterface
	userServiceref UserServiceInterface
}

func NewService(s TaskStoreInterface, us UserServiceInterface) *TaskService {
	return &TaskService{
		str:            s,
		userServiceref: us,
	}
}

func (s *TaskService) Create(t task.Task) (task.Task, error) {
	if err := t.Validate(); err != nil {
		return t, err
	}

	_, err := s.userServiceref.Get(t.Userid)
	if err != nil {
		return t, fmt.Errorf("user with ID %d does not exist: %v", t.Userid, err)
	}

	return s.str.CreateTask(t)
}

func (s *TaskService) GetTask(id int) (task.Task, error) {
	return s.str.GetByIDTask(id)
}

func (s *TaskService) Complete(id int) error {
	return s.str.CompleteTask(id)
}

func (s *TaskService) Delete(id int) error {
	return s.str.DeleteTask(id)
}

func (s *TaskService) All() ([]task.Task, error) {
	return s.str.GetAllTask()
}

func (s *TaskService) GetTasksByUserID(userid int) ([]task.Task, error) {
	_, err := s.userServiceref.Get(userid)

	if err != nil {
		return nil, err
	}

	return s.str.GetTasksByUserIDTask(userid)
}
