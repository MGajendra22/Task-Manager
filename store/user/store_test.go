package user

import (
	model "Task_Manager/model/user"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) (*UserStore, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return NewUserStore(db), mock, func() { _ = db.Close() }
}

func Test_CreateUser(t *testing.T) {
	store, mock, cleanup := setupDB(t)
	defer cleanup()

	u := model.User{Name: "John", Email: "john@example.com"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email) VALUES (?, ?)")).
		WithArgs(u.Name, u.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	created, err := store.CreateUser(u)
	require.NoError(t, err)
	require.Equal(t, 1, created.ID)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email) VALUES (?, ?)")).
		WithArgs(u.Name, u.Email).
		WillReturnError(errors.New("insert failed"))
	_, err = store.CreateUser(u)
	require.Error(t, err)
}

func Test_GetByIDUser(t *testing.T) {
	store, mock, cleanup := setupDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John", "john@example.com"))

	u, err := store.GetByIDUser(1)
	require.NoError(t, err)
	require.Equal(t, 1, u.ID)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users WHERE id = ?")).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)
	_, err = store.GetByIDUser(999)
	require.Error(t, err)
}

func Test_DeleteUser(t *testing.T) {
	store, mock, cleanup := setupDB(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	err := store.DeleteUser(1)
	require.NoError(t, err)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = ?")).
		WithArgs(999).
		WillReturnError(errors.New("delete failed"))
	err = store.DeleteUser(999)
	require.Error(t, err)
}

func Test_GetAllUser(t *testing.T) {
	store, mock, cleanup := setupDB(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John", "john@example.com").
			AddRow(2, "Alice", "alice@example.com")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users")).
			WillReturnRows(rows)

		users, err := store.GetAllUser()
		require.NoError(t, err)
		require.Len(t, users, 2)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users")).
			WillReturnError(errors.New("query failed"))

		_, err := store.GetAllUser()
		require.Error(t, err)
	})

	t.Run("scan error", func(t *testing.T) {
		// Intentionally missing one column to trigger scan error
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "John")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users")).
			WillReturnRows(rows)

		_, err := store.GetAllUser()
		require.Error(t, err)
	})
}
