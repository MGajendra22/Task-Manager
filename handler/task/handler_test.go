package task

import (
	"Task_Manager/model/task"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// errReader : Functionality is used to pass the empty and incorrect body to handle the edge case
type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

// errResponseWriter : Functionality is used to handle Write() method error
type errResponseWriter struct {
	http.ResponseWriter
}

func (e *errResponseWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("simulated write error")
}

// Test_NewHandler : To test that interface is correctly implemented or not
func Test_NewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockTaskServiceInterface(ctrl)

	h := NewHandler(mockSvc)

	if h == nil {
		t.Fatal("Expected non-nil handler")
	}

	if h.svc != mockSvc {
		t.Error("Expected service to be assigned correctly")
	}
}

// Test_CreateTasK : Tests task is created or not
func Test_CreateTasK(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		input       interface{}
		mockOutput  task.Task
		mockError   error
		ExpCode     int
		isWriteErr  bool
	}{
		{"Valid Input", "application/json", task.Task{1, "Working", false, 1}, task.Task{1, "Working", false, 1}, nil, http.StatusCreated, false},
		{"Invalid Json", "application/json", "bad json", task.Task{}, nil, http.StatusBadRequest, false},
		{"Creation error", "application/json", task.Task{1, "Working", false, 1}, task.Task{}, errors.New("creation error"), http.StatusBadRequest, false},
		{"wrong HTTP method", "application/json", task.Task{1, "Working", false, 1}, task.Task{}, nil, http.StatusMethodNotAllowed, false},
		{"Write Error", "application/json", task.Task{1, "Working", false, 1}, task.Task{1, "Working", false, 1}, nil, http.StatusCreated, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockTaskServiceInterface(ctrl)

			h := &Handler{mock}

			method := http.MethodPost
			if tt.name == "wrong HTTP method" {
				method = http.MethodPut
			}

			var body []byte
			if str, ok := tt.input.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			if method == http.MethodPost && tt.name != "Invalid Json" {
				mock.EXPECT().Create(gomock.Any()).Return(tt.mockOutput, tt.mockError).AnyTimes()
			}

			req := httptest.NewRequest(method, "/task", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)

			rec := httptest.NewRecorder()

			var w http.ResponseWriter = rec

			if tt.isWriteErr {
				w = &errResponseWriter{ResponseWriter: rec}
			}

			h.Create(w, req)

			if rec.Code != tt.ExpCode {
				t.Errorf("[%s] Expected status %d, got %d", tt.name, tt.ExpCode, rec.Code)
			}
		})
	}
}

// Test_CreateTask_ReadBodyError : To check the ReadBody Error
func Test_CreateTask_ReadBodyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := &Handler{NewMockTaskServiceInterface(ctrl)}

	r := httptest.NewRequest(http.MethodPost, "/task", errReader{})

	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test_GetTasK : Tests task is retrieved or not
func Test_GetTasK(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		ExpOutput  task.Task
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{"valid id", "1", task.Task{1, "Working", false, 1}, nil, http.StatusOK, false},
		{"Invalid user id", "abc", task.Task{1, "Working", false, 1}, nil, http.StatusBadRequest, false},
		{"Id not found", "99", task.Task{}, errors.New("Id Not found"), http.StatusNotFound, false},
		{"wrong HTTP method", "abc", task.Task{}, nil, http.StatusMethodNotAllowed, false},
		{"Write Error", "1", task.Task{1, "Working", false, 1}, nil, http.StatusOK, true},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := NewMockTaskServiceInterface(ctrl)

		h := &Handler{mock}

		method := http.MethodGet
		if tt.name == "wrong HTTP method" {
			method = http.MethodPost
		}

		if tt.ExpErr != nil || tt.ExpCode == http.StatusOK || tt.isWriteErr {
			mock.EXPECT().GetTask(gomock.Any()).Return(tt.ExpOutput, tt.ExpErr).AnyTimes()
		}

		req := httptest.NewRequest(method, "/task/"+tt.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": tt.id})

		rec := httptest.NewRecorder()

		var w http.ResponseWriter = rec

		if tt.isWriteErr {
			w = &errResponseWriter{ResponseWriter: rec}
		}

		h.GetTask(w, req)

		if rec.Code != tt.ExpCode {
			t.Errorf("Expected status %d, got %d", tt.ExpCode, rec.Code)
		}
	}
}

// Test_GetTasksByUserID : Tests task is retrieved or not of User
func Test_GetTasksByUserID(t *testing.T) {
	tests := []struct {
		name       string
		userid     string
		ExpOutput  []task.Task
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{"valid id", "1", []task.Task{{1, "Working", false, 1}}, nil, http.StatusOK, false},
		{"Invalid user id", "abc", nil, nil, http.StatusBadRequest, false},
		{"Id not found", "99", []task.Task{}, errors.New("Id Not found"), http.StatusNotFound, false},
		{"wrong HTTP method", "1", nil, nil, http.StatusMethodNotAllowed, false},
		{"Write Error", "1", []task.Task{{1, "Working", false, 1}}, nil, http.StatusOK, true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockTaskServiceInterface(ctrl)

	for _, tt := range tests {
		h := &Handler{mock}

		method := http.MethodGet
		if tt.name == "wrong HTTP method" {
			method = http.MethodPost
		}

		if id, err := strconv.Atoi(tt.userid); err == nil && method == http.MethodGet && (tt.ExpCode != http.StatusBadRequest || tt.isWriteErr) {
			mock.EXPECT().GetTasksByUserID(id).Return(tt.ExpOutput, tt.ExpErr).AnyTimes()
		}

		req := httptest.NewRequest(method, "/task/user/"+tt.userid, nil)
		req = mux.SetURLVars(req, map[string]string{"userid": tt.userid})
		rec := httptest.NewRecorder()

		var w http.ResponseWriter = rec

		if tt.isWriteErr {
			w = &errResponseWriter{ResponseWriter: rec}
		}

		h.GetTasksByUserID(w, req)

		if rec.Code != tt.ExpCode {
			t.Errorf("Expected status %d, got %d", tt.ExpCode, rec.Code)
		}
	}
}

// Test_CompleteTask : Tests assigned id's task is completed or not
func Test_CompleteTask(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		ExpOutput  task.Task
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{"valid id", "1", task.Task{1, "Working", true, 1}, nil, http.StatusOK, false},
		{"Invalid user id", "abc", task.Task{1, "Working", true, 1}, nil, http.StatusBadRequest, false},
		{"Id not found", "99", task.Task{}, errors.New("Id Not found"), http.StatusNotFound, false},
		{"wrong HTTP method", "1", task.Task{1, "Working", true, 1}, nil, http.StatusMethodNotAllowed, false},
		{"Write Error", "1", task.Task{1, "Working", true, 1}, nil, http.StatusOK, true},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockTaskServiceInterface(ctrl)

		h := &Handler{mock}

		method := http.MethodPut
		if tt.name == "wrong HTTP method" {
			method = http.MethodGet
		}

		if id, err := strconv.Atoi(tt.id); err == nil && method == http.MethodPut && (tt.ExpCode != http.StatusBadRequest || tt.isWriteErr) {
			mock.EXPECT().Complete(id).Return(tt.ExpErr).AnyTimes()
		}

		req := httptest.NewRequest(method, "/task/"+tt.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": tt.id})

		rec := httptest.NewRecorder()

		var w http.ResponseWriter = rec

		if tt.isWriteErr {
			w = &errResponseWriter{ResponseWriter: rec}
		}

		h.Complete(w, req)

		if rec.Code != tt.ExpCode {
			t.Errorf("Expected status %d, got %d", tt.ExpCode, rec.Code)
		}
	}

}

// Test_AllTasks : Tests all tasks are retrieved or not
func Test_AllTasks(t *testing.T) {
	tests := []struct {
		name       string
		input      []task.Task
		ExpOutput  []task.Task
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{
			name:      "Successfully retrieved",
			input:     []task.Task{{1, "Working", true, 1}},
			ExpOutput: []task.Task{{1, "Working", true, 1}},
			ExpErr:    nil,
			ExpCode:   http.StatusOK,
		},
		{
			name:      "Unable to fetch user data",
			input:     nil,
			ExpOutput: nil,
			ExpErr:    errors.New("Failed to fetch user's data"),
			ExpCode:   http.StatusInternalServerError,
		},
		{
			name:      "wrong HTTP method",
			input:     nil,
			ExpOutput: nil,
			ExpErr:    nil,
			ExpCode:   http.StatusMethodNotAllowed,
		},
		{
			name:       "Write failure after successful fetch",
			input:      []task.Task{{1, "desc", true, 1}},
			ExpOutput:  []task.Task{{1, "desc", true, 1}},
			ExpErr:     nil,
			ExpCode:    http.StatusOK, // write failure still results in 200 OK
			isWriteErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockTaskServiceInterface(ctrl)

			h := &Handler{mock}

			method := http.MethodGet
			if tt.name == "wrong HTTP method" {
				method = http.MethodPost
			}

			if tt.ExpErr != nil || tt.ExpCode == http.StatusOK || tt.isWriteErr {
				mock.EXPECT().All().Return(tt.ExpOutput, tt.ExpErr).AnyTimes()
			}

			req := httptest.NewRequest(method, "/task", nil)

			rec := httptest.NewRecorder()

			var w http.ResponseWriter = rec

			if tt.isWriteErr {
				w = &errResponseWriter{ResponseWriter: rec}
			}

			h.All(w, req)

			if rec.Code != tt.ExpCode {
				t.Errorf("[%s] Expected status %d, got %d", tt.name, tt.ExpCode, rec.Code)
			}
		})
	}
}

// Test_DeleteTask : Tests Task with task-id is deleted or not
func Test_DeleteTask(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{"Valid delete", "1", nil, http.StatusOK, false},
		{"InValid user id", "abc", errors.New("Invalid user id "), http.StatusBadRequest, false},
		{"Delete error", "99", errors.New("Delete error"), http.StatusNotFound, false},
		{"wrong HTTP method", "abc", nil, http.StatusMethodNotAllowed, false},
		{"Write Error", "1", nil, http.StatusOK, true},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := NewMockTaskServiceInterface(ctrl)

		if tt.ExpErr != nil || tt.ExpCode == http.StatusOK || tt.isWriteErr {
			mock.EXPECT().Delete(gomock.Any()).Return(tt.ExpErr).AnyTimes()
		}

		method := http.MethodDelete
		if tt.name == "wrong HTTP method" {
			method = http.MethodGet
		}

		h := &Handler{mock}

		req := httptest.NewRequest(method, "/task/"+tt.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": tt.id})

		rec := httptest.NewRecorder()

		var w http.ResponseWriter = rec
		if tt.isWriteErr {
			w = &errResponseWriter{ResponseWriter: rec}
		}

		h.Delete(w, req)

		if rec.Code != tt.ExpCode {
			t.Errorf("Expected status %d, got %d", tt.ExpCode, rec.Code)
		}
	}
}
