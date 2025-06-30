package user

import (
	"Task_Manager/model/user"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type mockUserService struct{}

func (m mockUserService) Create(u user.User) (user.User, error) {
	if u.Name == "fail" {
		return user.User{}, errors.New("create failed")
	}
	u.ID = 1
	return u, nil
}

func (m mockUserService) Get(id int) (user.User, error) {
	if id == 0 {
		return user.User{}, errors.New("not found")
	}
	return user.User{ID: id, Name: "Test"}, nil
}

func (m mockUserService) Delete(id int) error {
	if id == 0 {
		return errors.New("delete failed")
	}
	return nil
}

func (m mockUserService) All() ([]user.User, error) {
	return []user.User{{ID: 1, Name: "Alice"}}, nil
}

func Test_CreateUser(t *testing.T) {
	h := NewUserHandler(mockUserService{})

	t.Run("invalid content-type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user", nil)
		req.Header.Set("Content-Type", "text/plain")
		rr := httptest.NewRecorder()
		h.CreateUser(rr, req)
		require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateUser(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("creation failure", func(t *testing.T) {
		body, _ := json.Marshal(user.User{Name: "fail"})
		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateUser(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("success", func(t *testing.T) {
		body, _ := json.Marshal(user.User{Name: "pass"})
		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateUser(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)
	})
}

func Test_GetUser(t *testing.T) {
	h := NewUserHandler(mockUserService{})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/x", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "x"})
		rr := httptest.NewRecorder()
		h.GetUser(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		rr := httptest.NewRecorder()
		h.GetUser(rr, req)
		require.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rr := httptest.NewRecorder()
		h.GetUser(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})
}

func Test_DeleteUser(t *testing.T) {
	h := NewUserHandler(mockUserService{})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/user/x", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "x"})
		rr := httptest.NewRecorder()
		h.DeleteUser(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("delete fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/user/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		rr := httptest.NewRecorder()
		h.DeleteUser(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rr := httptest.NewRecorder()
		h.DeleteUser(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})
}

func Test_GetAllUsers(t *testing.T) {
	h := NewUserHandler(mockUserService{})
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rr := httptest.NewRecorder()
	h.GetAllUsers(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	body, _ := io.ReadAll(rr.Body)
	require.Contains(t, string(body), "Alice")
}
