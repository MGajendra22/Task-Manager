package user

import (
	"Task_Manager/model/user"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
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

// Test_UserHandler : To test that interface is correctly implemented or not
func Test_UserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(mockSvc)

	if h == nil {
		t.Fatal("Expected non-nil handler")
	}
	if h.Service != mockSvc {
		t.Error("Expected service to be assigned correctly")
	}
}

// Test_CreateUseR : Tests user is task is created or not
func Test_CreateUseR(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		input        interface{}
		mockReturn   user.User
		mockError    error
		expectedCode int
	}{
		{
			name:         "valid input",
			contentType:  "application/json",
			input:        user.User{Name: "John", Email: "john@gmail.com"},
			mockReturn:   user.User{ID: 1, Name: "John", Email: "john@gmail.com"},
			mockError:    nil,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "creation failure",
			contentType:  "application/json",
			input:        user.User{Name: "john", Email: "john@email"},
			mockReturn:   user.User{},
			mockError:    errors.New("creation error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "bad JSON",
			contentType:  "application/json",
			input:        "{bad json}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "wrong HTTP method",
			contentType:  "application/json",
			input:        user.User{Name: "Wrong", Email: "wrong@example.com"},
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockUserServiceInterface(ctrl)
			handler := &UserHandler{Service: mockService}

			method := http.MethodPost
			if tt.name == "wrong HTTP method" {
				method = http.MethodGet
			}

			var body []byte
			if str, ok := tt.input.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest(method, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			if tt.expectedCode == http.StatusCreated || tt.expectedCode == http.StatusInternalServerError {
				mockService.EXPECT().Create(gomock.Any()).Return(tt.mockReturn, tt.mockError).AnyTimes()
			}

			handler.CreateUser(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

// Test_CreateUser_ReadBodyError : To check the ReadBody Error
func Test_CreateUser_ReadBodyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := &UserHandler{Service: NewMockUserServiceInterface(ctrl)}

	r := httptest.NewRequest(http.MethodPost, "/users", errReader{})
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test_GetUseR : Tests user details are retrieved or not
func Test_GetUseR(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		ExpOutput  user.User
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{"valid id", "1", user.User{ID: 1, Name: "john", Email: "john@gamil.com"}, nil, http.StatusOK, false},
		{"Invalid user id", "abc", user.User{ID: 1, Name: "john", Email: "john@gamil.com"}, nil, http.StatusBadRequest, false},
		{"Id not found", "99", user.User{}, errors.New("Id Not found"), http.StatusNotFound, false},
		{"wrong HTTP method", "abc", user.User{}, nil, http.StatusMethodNotAllowed, false},
		{"Write Error", "1", user.User{ID: 1, Name: "john", Email: "john@gamil.com"}, nil, http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockUserServiceInterface(ctrl)

			h := &UserHandler{mock}
			method := http.MethodGet
			if tt.name == "wrong HTTP method" {
				method = http.MethodPost
			}

			if tt.ExpErr != nil || tt.ExpCode == http.StatusOK || tt.isWriteErr {
				mock.EXPECT().Get(gomock.Any()).Return(tt.ExpOutput, tt.ExpErr).AnyTimes()
			}

			req := httptest.NewRequest(method, "/users/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rec := httptest.NewRecorder()
			var w http.ResponseWriter = rec
			if tt.isWriteErr {
				w = &errResponseWriter{ResponseWriter: rec}
			}
			h.GetUser(w, req)

			if rec.Code != tt.ExpCode {
				t.Errorf("GetUser1() = %v, want %v", rec.Code, tt.ExpCode)
			}
		})
	}
}

// Test_DeleteUseR : Tests user with user-id is deleted or not
func Test_DeleteUseR(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{"Valid delete", "1", nil, http.StatusOK, false},
		{"InValid user id", "abc", errors.New("Invalid user id "), http.StatusBadRequest, false},
		{"Delete error", "99", errors.New("Delete error"), http.StatusInternalServerError, false},
		{"wrong HTTP method", "abc", nil, http.StatusMethodNotAllowed, false},
		{"Write Error", "1", nil, http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockUserServiceInterface(ctrl)
			h := &UserHandler{mock}

			method := http.MethodDelete
			if tt.name == "wrong HTTP method" {
				method = http.MethodGet
			}
			if tt.ExpErr != nil || tt.ExpCode == http.StatusOK || tt.isWriteErr {
				mock.EXPECT().Delete(gomock.Any()).Return(tt.ExpErr).AnyTimes()
			}

			req := httptest.NewRequest(method, "/users/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rec := httptest.NewRecorder()
			var w http.ResponseWriter = rec
			if tt.isWriteErr {
				w = &errResponseWriter{ResponseWriter: rec}
			}
			h.DeleteUser(w, req)

			if rec.Code != tt.ExpCode {
				t.Errorf("DeleteUser1() = %v, want %v", rec.Code, tt.ExpCode)
			}
		})
	}
}

// Test_GetAllUsers : Tests all user details are retrieved or not
func Test_GetAllUsers(t *testing.T) {
	tests := []struct {
		name       string
		input      []user.User
		ExpOutput  []user.User
		ExpErr     error
		ExpCode    int
		isWriteErr bool
	}{
		{"Successfully retried", []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, nil, http.StatusOK, false},
		{"Unable to fetch user data", []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, nil, errors.New("Failed to fetch user's data"), http.StatusInternalServerError, false},
		{"wrong HTTP method", []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, nil, nil, http.StatusMethodNotAllowed, false},
		{"Write Error", []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, nil, http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockUserServiceInterface(ctrl)
			h := &UserHandler{mock}

			method := http.MethodGet
			if tt.name == "wrong HTTP method" {
				method = http.MethodPost
			}
			if tt.ExpErr != nil || tt.ExpCode == http.StatusOK || tt.isWriteErr {
				mock.EXPECT().All().Return(tt.ExpOutput, tt.ExpErr).AnyTimes()
			}

			req := httptest.NewRequest(method, "/users", nil)
			rec := httptest.NewRecorder()
			var w http.ResponseWriter = rec
			if tt.isWriteErr {
				w = &errResponseWriter{ResponseWriter: rec}
			}
			h.GetAllUsers(w, req)
			if rec.Code != tt.ExpCode {
				t.Errorf("GetAllUsers1() = %v, want %v", rec.Code, tt.ExpCode)
			}

		})
	}
}
