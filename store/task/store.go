package task

import (
	"Task_Manager/model/task"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateTask inserts a new task into the database
func (s *Store) CreateTask(t task.Task) (task.Task, error) {
	res, err := s.db.Exec("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)", t.Desc, t.Status, t.Userid)

	if err != nil {
		return t, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return t, err
	}

	t.ID = int(id)

	return t, nil
}

// GetByIDTask fetches a task by its ID
func (s *Store) GetByIDTask(id int) (task.Task, error) {
	var t task.Task
	err := s.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Desc, &t.Status, &t.Userid)

	if err != nil {
		return t, err
	}

	return t, err
}

// CompleteTask marks a task as completed
func (s *Store) CompleteTask(id int) error {
	res, err := s.db.Exec("UPDATE tasks SET status = true WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteTask removes a task by ID
func (s *Store) DeleteTask(id int) error {
	res, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetAllTask returns all tasks from the database
func (s *Store) GetAllTask() ([]task.Task, error) {
	rows, err := s.db.Query("SELECT id, description, status , userid FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.Userid); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTasksByUserID it will send the tasks , which are assigned to user
func (s *Store) GetTasksByUserIDTask(userid int) ([]task.Task, error) {
	rows, err := s.db.Query("SELECT id, description, status , userid FROM tasks where userid =?", userid)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var tasks []task.Task

	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.Userid); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
