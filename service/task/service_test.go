package task

import (
	taskModel "Task_Manager/model/task"
	"Task_Manager/model/user"
	"errors"
	"testing"
)

// Mock TaskStoreInterface
type mockTaskStore struct{}

func (m mockTaskStore) CreateTask(t taskModel.Task) (taskModel.Task, error) {
	t.ID = 101
	return t, nil
}
func (m mockTaskStore) GetByIDTask(id int) (taskModel.Task, error) {
	if id == 0 {
		return taskModel.Task{}, errors.New("not found")
	}
	return taskModel.Task{ID: id, Desc: "Test", Userid: 1}, nil
}
func (m mockTaskStore) GetAllTask() ([]taskModel.Task, error) {
	return []taskModel.Task{{ID: 1, Desc: "Task1"}}, nil
}
func (m mockTaskStore) CompleteTask(id int) error {
	if id == 999 {
		return errors.New("not found")
	}
	return nil
}
func (m mockTaskStore) DeleteTask(id int) error {
	if id == 999 {
		return errors.New("not found")
	}
	return nil
}
func (m mockTaskStore) GetTasksByUserIDTask(userId int) ([]taskModel.Task, error) {
	if userId <= 0 {
		return nil, errors.New("invalid user")
	}
	return []taskModel.Task{{ID: 1, Desc: "UserTask", Userid: userId}}, nil
}

// Mock UserServiceInterface
type mockUserService struct{}

func (m mockUserService) Get(id int) (user.User, error) {
	if id <= 0 {
		return user.User{}, errors.New("user not found")
	}
	return user.User{ID: id, Name: "John"}, nil
}

func TestService_Create(t *testing.T) {
	svc := NewService(mockTaskStore{}, mockUserService{})

	t.Run("invalid task (validation fails)", func(t *testing.T) {
		invalid := taskModel.Task{Desc: "", Userid: 1}
		_, err := svc.Create(invalid)
		if err == nil {
			t.Errorf("expected validation error")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		input := taskModel.Task{Desc: "Valid", Userid: -1}
		_, err := svc.Create(input)
		if err == nil {
			t.Errorf("expected user not found error")
		}
	})

	t.Run("success", func(t *testing.T) {
		input := taskModel.Task{Desc: "Do this", Userid: 1}
		createdTask, err := svc.Create(input)
		if err != nil || createdTask.ID != 101 {
			t.Errorf("expected task ID 101, got %v, err: %v", createdTask.ID, err)
		}
	})
}

func TestService_GetTask(t *testing.T) {
	svc := NewService(mockTaskStore{}, mockUserService{})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetTask(0)
		if err == nil {
			t.Errorf("expected error for id=0")
		}
	})

	t.Run("success", func(t *testing.T) {
		result, err := svc.GetTask(5)
		if err != nil || result.ID != 5 {
			t.Errorf("unexpected task result: %+v, err: %v", result, err)
		}
	})
}

func TestService_Complete(t *testing.T) {
	svc := NewService(mockTaskStore{}, mockUserService{})

	t.Run("invalid ID", func(t *testing.T) {
		err := svc.Complete(999)
		if err == nil {
			t.Errorf("expected error for invalid id")
		}
	})

	t.Run("success", func(t *testing.T) {
		err := svc.Complete(1)
		if err != nil {
			t.Errorf("expected success, got error: %v", err)
		}
	})
}

func TestService_Delete(t *testing.T) {
	svc := NewService(mockTaskStore{}, mockUserService{})

	t.Run("invalid ID", func(t *testing.T) {
		err := svc.Delete(999)
		if err == nil {
			t.Errorf("expected error for invalid id")
		}
	})

	t.Run("success", func(t *testing.T) {
		err := svc.Delete(1)
		if err != nil {
			t.Errorf("expected success, got error: %v", err)
		}
	})
}

func TestService_All(t *testing.T) {
	svc := NewService(mockTaskStore{}, mockUserService{})

	t.Run("returns all tasks", func(t *testing.T) {
		tasks, err := svc.All()
		if err != nil || len(tasks) != 1 {
			t.Errorf("expected 1 task, got %v, err: %v", tasks, err)
		}
	})
}

func TestService_GetTasksByUserID(t *testing.T) {
	svc := NewService(mockTaskStore{}, mockUserService{})

	t.Run("user not found", func(t *testing.T) {
		_, err := svc.GetTasksByUserID(-1)
		if err == nil {
			t.Errorf("expected error for invalid user")
		}
	})

	t.Run("success", func(t *testing.T) {
		tasks, err := svc.GetTasksByUserID(1)
		if err != nil || len(tasks) != 1 || tasks[0].Userid != 1 {
			t.Errorf("unexpected tasks: %+v, err: %v", tasks, err)
		}
	})
}
