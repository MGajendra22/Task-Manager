package user

import (
	"Task_Manager/model/user"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockUserStore struct{}

func (m mockUserStore) CreateUser(u user.User) (user.User, error) {
	if u.Name == "fail" {
		return user.User{}, errors.New("create failed")
	}
	u.ID = 1
	return u, nil
}

func (m mockUserStore) GetByIDUser(id int) (user.User, error) {
	if id <= 0 {
		return user.User{}, errors.New("not found")
	}
	return user.User{ID: id, Name: "John"}, nil
}

func (m mockUserStore) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("delete error")
	}
	return nil
}

func (m mockUserStore) GetAllUser() ([]user.User, error) {
	return []user.User{{ID: 1, Name: "A"}}, nil
}

func TestUserService_Create(t *testing.T) {
	svc := NewUserService(mockUserStore{})

	t.Run("validation error", func(t *testing.T) {
		u := user.User{Name: ""}
		_, err := svc.Create(u)
		require.Error(t, err)
	})

	t.Run("store error", func(t *testing.T) {
		u := user.User{Name: "fail", Email: "x@y.com"}
		_, err := svc.Create(u)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		u := user.User{Name: "John", Email: "a@b.com"}
		created, err := svc.Create(u)
		require.NoError(t, err)
		require.Equal(t, 1, created.ID)
	})
}

func TestUserService_Get(t *testing.T) {
	svc := NewUserService(mockUserStore{})
	_, err := svc.Get(0)
	require.Error(t, err)
	usr, err := svc.Get(2)
	require.NoError(t, err)
	require.Equal(t, 2, usr.ID)
}

func TestUserService_Delete(t *testing.T) {
	svc := NewUserService(mockUserStore{})
	err := svc.Delete(0)
	require.Error(t, err)
	err = svc.Delete(1)
	require.NoError(t, err)
}

func TestUserService_All(t *testing.T) {
	svc := NewUserService(mockUserStore{})
	users, err := svc.All()
	require.NoError(t, err)
	require.Len(t, users, 1)
}
