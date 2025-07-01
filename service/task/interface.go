package task

import (
	"Task_Manager/model/task"
	userModel "Task_Manager/model/user"
)

type TaskStoreInterface interface {
	CreateTask(task task.Task) (task.Task, error)
	GetByIDTask(id int) (task.Task, error)
	GetAllTask() ([]task.Task, error)
	CompleteTask(id int) error
	DeleteTask(id int) error
	GetTasksByUserIDTask(userId int) ([]task.Task, error)
}

type UserServiceInterface interface {
	Get(id int) (userModel.User, error)
}
