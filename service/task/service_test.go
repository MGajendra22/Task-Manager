package task

import (
	"Task_Manager/model/task"
	"Task_Manager/model/user"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_Create(t *testing.T) {
	tests := []struct {
		name        string
		input       task.Task
		mockUser    user.User
		mockTaskOut task.Task
		userErr     error
		taskErr     error
		expErr      bool
	}{
		{
			name:        "Valid Task Creation",
			input:       task.Task{ID: 1, Desc: "Do Work", Status: false, Userid: 10},
			mockUser:    user.User{ID: 10, Name: "Alice", Email: "alice@example.com"},
			mockTaskOut: task.Task{ID: 1, Desc: "Do Work", Status: false, Userid: 10},
			expErr:      false,
		},
		{
			name:   "Validation Error - Empty Desc",
			input:  task.Task{ID: 2, Desc: "", Userid: 10},
			expErr: true,
		},
		{
			name:    "User Not Found",
			input:   task.Task{ID: 3, Desc: "Plan", Userid: 20},
			userErr: errors.New("user not found"),
			expErr:  true,
		},
		{
			name:     "Task Store Error",
			input:    task.Task{ID: 4, Desc: "Build", Userid: 11},
			mockUser: user.User{ID: 11, Name: "Bob", Email: "bob@example.com"},
			taskErr:  errors.New("db write failed"),
			expErr:   true,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := NewMockTaskStoreInterface(ctrl)
		mockUserServ := NewMockUserServiceInterface(ctrl)
		service := NewService(mockStore, mockUserServ)

		if err := tt.input.Validate(); err == nil {
			mockUserServ.EXPECT().
				Get(tt.input.Userid).
				Return(tt.mockUser, tt.userErr)

			if tt.userErr == nil {
				mockStore.EXPECT().
					CreateTask(tt.input).
					Return(tt.mockTaskOut, tt.taskErr)
			}
		}

		result, err := service.Create(tt.input)

		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.Equal(t, tt.mockTaskOut, result, tt.name)
		}
	}
}

func Test_GetTask(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		mockOutput task.Task
		mockErr    error
		expErr     bool
	}{
		{"Valid Id", 1, task.Task{1, "Working", false, 1}, nil, false},
		{"Task Not found", 1, task.Task{}, errors.New("task not found"), true},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStore := NewMockTaskStoreInterface(ctrl)
		mockUserServ := NewMockUserServiceInterface(ctrl)
		service := NewService(mockStore, mockUserServ)

		mockStore.EXPECT().
			GetByIDTask(tt.id).
			Return(tt.mockOutput, tt.mockErr)

		res, err := service.GetTask(tt.id)

		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.Equal(t, tt.mockOutput, res, tt.name)
		}

	}
}

func Test_AllTasks(t *testing.T) {
	tests := []struct {
		name       string
		mockOutput []task.Task
		mockErr    error
		expErr     bool
	}{
		{"Data fetched", []task.Task{{1, "Working", false, 1}}, nil, false},
		{"Unable to fetch", []task.Task{}, errors.New("task not found"), true},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStore := NewMockTaskStoreInterface(ctrl)
		mockUserServ := NewMockUserServiceInterface(ctrl)
		service := NewService(mockStore, mockUserServ)
		mockStore.EXPECT().GetAllTask().Return(tt.mockOutput, tt.mockErr)

		res, err := service.All()
		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.ElementsMatch(t, tt.mockOutput, res, tt.name)
		}

	}
}

func Test_CompleteTask(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		taskErr error
		expErr  bool
	}{
		{
			name:   "Valid Task Creation",
			input:  1,
			expErr: false,
		},
		{
			name:    "Task Not Found",
			input:   1,
			taskErr: errors.New("user not found"),
			expErr:  true,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStore := NewMockTaskStoreInterface(ctrl)

		service := NewService(mockStore, nil)
		mockStore.EXPECT().CompleteTask(tt.input).Return(tt.taskErr).AnyTimes()

		err := service.Complete(tt.input)
		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)

		}

	}
}

func Test_DeleteTask(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		taskErr error
		expErr  bool
	}{
		{
			name:   "Valid Task Deletion",
			input:  1,
			expErr: false,
		},
		{
			name:    "Task Not Found",
			input:   1,
			taskErr: errors.New("Task not found"),
			expErr:  true,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStore := NewMockTaskStoreInterface(ctrl)
		service := NewService(mockStore, nil)
		mockStore.EXPECT().DeleteTask(tt.input).Return(tt.taskErr).AnyTimes()
		err := service.Delete(tt.input)
		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
	}

}

func Test_GetTasksByUserId(t *testing.T) {
	tests := []struct {
		name       string
		input      int
		mockUser   user.User
		userErr    error
		mockOutput []task.Task
		mockErr    error
		expErr     bool
	}{
		{
			name:       "Data fetched",
			input:      1,
			mockUser:   user.User{ID: 1, Name: "Test", Email: "test@test.com"},
			userErr:    nil,
			mockOutput: []task.Task{{ID: 1, Desc: "Working", Status: false, Userid: 1}},
			mockErr:    nil,
			expErr:     false,
		},
		{
			name:       "User not found",
			input:      2,
			userErr:    errors.New("User Not found"),
			mockOutput: nil,
			expErr:     true,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := NewMockTaskStoreInterface(ctrl)
		mockUserServ := NewMockUserServiceInterface(ctrl)
		service := NewService(mockStore, mockUserServ)

		mockUserServ.EXPECT().
			Get(tt.input).
			Return(tt.mockUser, tt.userErr)

		if tt.userErr == nil {
			mockStore.EXPECT().
				GetTasksByUserIDTask(tt.input).
				Return(tt.mockOutput, tt.mockErr)
		}

		res, err := service.GetTasksByUserID(tt.input)

		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.ElementsMatch(t, tt.mockOutput, res, tt.name)
		}
	}
}
