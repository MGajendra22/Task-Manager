package task

import (
	"Task_Manager/model/task"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// 👇 Corrected error-inducing io.ReadCloser mock
type errReadCloser struct{}

func (errReadCloser) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (errReadCloser) Close() error             { return nil }

// ✅ Mock TaskService
type mockTaskService struct{}

func (m mockTaskService) Create(t task.Task) (task.Task, error) {
	if t.Desc == "" {
		return task.Task{}, errors.New("desc empty")
	}
	t.ID = 42
	return t, nil
}

func (m mockTaskService) GetTask(id int) (task.Task, error) {
	if id <= 0 {
		return task.Task{}, errors.New("not found")
	}
	return task.Task{ID: id, Desc: "Loaded", Userid: 1}, nil
}

func (m mockTaskService) GetTasksByUserID(userId int) ([]task.Task, error) {
	if userId <= 0 {
		return nil, errors.New("task not found")
	}
	return []task.Task{{ID: 1, Desc: "ByUser", Userid: userId}}, nil
}

func (m mockTaskService) Complete(id int) error {
	if id != 1 {
		return errors.New("not found")
	}
	return nil
}

func (m mockTaskService) Delete(id int) error {
	if id != 1 {
		return errors.New("not found")
	}
	return nil
}

func (m mockTaskService) All() ([]task.Task, error) {
	return []task.Task{{ID: 1, Desc: "AllTask", Userid: 1}}, nil
}

// 🔁 Reusable helper to create requests with path variables
func newReq(method, url, body string, vars map[string]string) (*http.Request, *httptest.ResponseRecorder) {
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, url, rdr)
	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}
	return req, httptest.NewRecorder()
}

func Test_Handler_Create(t *testing.T) {
	h := NewHandler(mockTaskService{})

	t.Run("wrong method", func(t *testing.T) {

		req, rr := newReq(http.MethodGet, "/task", "", nil)
		h.Create(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rr.Code)
		}
	})

	t.Run("read error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/task", nil)
		req.Body = errReadCloser{}
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req, rr := newReq(http.MethodPost, "/task", "{bad}", nil)
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		req, rr := newReq(http.MethodPost, "/task", `{"desc":"","userid":1}`, nil)
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		req, rr := newReq(http.MethodPost, "/task", `{"desc":"Hello","userid":1}`, nil)
		h.Create(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rr.Code)
		}
	})
}

func Test_Handler_GetTask(t *testing.T) {
	h := NewHandler(mockTaskService{})

	t.Run("wrong method", func(t *testing.T) {
		req, rr := newReq(http.MethodPost, "/task/1", "", map[string]string{"id": "1"})
		h.GetTask(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rr.Code)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, rr := newReq(http.MethodGet, "/task/bad", "", map[string]string{"id": "bad"})
		h.GetTask(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req, rr := newReq(http.MethodGet, "/task/0", "", map[string]string{"id": "0"})
		h.GetTask(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		req, rr := newReq(http.MethodGet, "/task/3", "", map[string]string{"id": "3"})
		h.GetTask(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	})
}

func Test_Handler_GetTasksByUserID(t *testing.T) {
	h := NewHandler(mockTaskService{})

	t.Run("wrong method", func(t *testing.T) {
		req, rr := newReq(http.MethodPost, "/task/user/1", "", map[string]string{"userid": "1"})
		h.GetTasksByUserID(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rr.Code)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, rr := newReq(http.MethodGet, "/task/user/bad", "", map[string]string{"userid": "bad"})
		h.GetTasksByUserID(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req, rr := newReq(http.MethodGet, "/task/user/0", "", map[string]string{"userid": "0"})
		h.GetTasksByUserID(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		req, rr := newReq(http.MethodGet, "/task/user/1", "", map[string]string{"userid": "1"})
		h.GetTasksByUserID(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	})
}

func Test_Handler_CompleteAndDelete(t *testing.T) {
	h := NewHandler(mockTaskService{})

	tests := []struct {
		name     string
		handler  func(http.ResponseWriter, *http.Request)
		method   string
		pathVar  string
		expected int
	}{
		{"Complete wrong method", h.Complete, http.MethodGet, "1", http.StatusMethodNotAllowed},
		{"Complete bad ID", h.Complete, http.MethodPut, "bad", http.StatusBadRequest},
		{"Complete not found", h.Complete, http.MethodPut, "2", http.StatusNotFound},
		{"Complete success", h.Complete, http.MethodPut, "1", http.StatusOK},

		{"Delete wrong method", h.Delete, http.MethodGet, "1", http.StatusMethodNotAllowed},
		{"Delete bad ID", h.Delete, http.MethodDelete, "bad", http.StatusBadRequest},
		{"Delete not found", h.Delete, http.MethodDelete, "2", http.StatusNotFound},
		{"Delete success", h.Delete, http.MethodDelete, "1", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, rr := newReq(tt.method, "/task", "", map[string]string{"id": tt.pathVar})
			tt.handler(rr, req)
			if rr.Code != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, rr.Code)
			}
		})
	}
}

func Test_Handler_All(t *testing.T) {
	h := NewHandler(mockTaskService{})

	t.Run("wrong method", func(t *testing.T) {
		req, rr := newReq(http.MethodPost, "/task", "", nil)
		h.All(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		req, rr := newReq(http.MethodGet, "/task", "", nil)
		h.All(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}

		var tasks []task.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil || len(tasks) != 1 {
			t.Errorf("expected 1 task, got %d and err %v", len(tasks), err)
		}
	})
}
