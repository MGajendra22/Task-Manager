package task

import "errors"

type Task struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status bool   `json:"status"`
	Userid int    `json:"userid"`
}

func (t *Task) Validate() error {
	if t.Desc == "" {
		return errors.New("description cannot be empty")
	}
	return nil
}
