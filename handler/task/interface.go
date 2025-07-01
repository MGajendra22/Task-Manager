package task

import "Task_Manager/model/task"

type TaskServiceInterface interface {
	Create(t task.Task) (task.Task, error)
	GetTask(id int) (task.Task, error)
	Complete(id int) error
	Delete(id int) error
	All() ([]task.Task, error)
	GetTasksByUserID(userId int) ([]task.Task, error)
}
