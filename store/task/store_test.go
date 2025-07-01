package task

import (
	taskModel "Task_Manager/model/task"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*Store, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	cleanup := func() { _ = db.Close() }
	return NewStore(db), mock, cleanup
}

func Test_CreateTask(t *testing.T) {
	store, mock, cleanup := setup(t)
	defer cleanup()

	tsk := taskModel.Task{Desc: "New Task", Status: false, Userid: 2}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)")).
			WithArgs(tsk.Desc, tsk.Status, tsk.Userid).
			WillReturnResult(sqlmock.NewResult(1, 1))

		created, err := store.CreateTask(tsk)
		require.NoError(t, err)
		require.Equal(t, 1, created.ID)
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)")).
			WithArgs(tsk.Desc, tsk.Status, tsk.Userid).
			WillReturnError(errors.New("insert failed"))

		_, err := store.CreateTask(tsk)
		require.Error(t, err)
		require.EqualError(t, err, "insert failed")
	})

	t.Run("LastInsertId Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)")).
			WithArgs(tsk.Desc, tsk.Status, tsk.Userid).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("lastInsertId failed")))

		_, err := store.CreateTask(tsk)
		require.Error(t, err)
		require.EqualError(t, err, "lastInsertId failed")
	})
}

func Test_GetByIDTask(t *testing.T) {
	store, mock, cleanup := setup(t)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM tasks WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "status", "userid"}).
				AddRow(1, "Do homework", false, 1))
		tsk, err := store.GetByIDTask(1)
		require.NoError(t, err)
		require.Equal(t, 1, tsk.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM tasks WHERE id = ?")).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)
		_, err := store.GetByIDTask(999)
		require.Error(t, err)
	})
}

func Test_CompleteTask(t *testing.T) {
	store, mock, cleanup := setup(t)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET status = true WHERE id = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		err := store.CompleteTask(1)
		require.NoError(t, err)
	})

	t.Run("No Rows Updated", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET status = true WHERE id = ?")).
			WithArgs(2).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := store.CompleteTask(2)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET status = true WHERE id = ?")).
			WithArgs(3).
			WillReturnError(errors.New("db error"))
		err := store.CompleteTask(3)
		require.Error(t, err)
	})

	t.Run("RowsAffected Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET status = true WHERE id = ?")).
			WithArgs(4).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected fail")))
		err := store.CompleteTask(4)
		require.Error(t, err)
	})
}

func Test_DeleteTask(t *testing.T) {
	store, mock, cleanup := setup(t)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM tasks WHERE id = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		err := store.DeleteTask(1)
		require.NoError(t, err)
	})

	t.Run("No Rows Deleted", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM tasks WHERE id = ?")).
			WithArgs(2).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := store.DeleteTask(2)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM tasks WHERE id = ?")).
			WithArgs(3).
			WillReturnError(sql.ErrConnDone)
		err := store.DeleteTask(3)
		require.Error(t, err)
	})

	t.Run("RowsAffected Error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM tasks WHERE id = ?")).
			WithArgs(4).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected fail")))
		err := store.DeleteTask(4)
		require.Error(t, err)
	})
}

func Test_GetAllTask(t *testing.T) {
	store, mock, cleanup := setup(t)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, status , userid FROM tasks")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "status", "userid"}).
				AddRow(1, "Task1", false, 1).
				AddRow(2, "Task2", true, 2))
		tasks, err := store.GetAllTask()
		require.NoError(t, err)
		require.Len(t, tasks, 2)
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, status , userid FROM tasks")).
			WillReturnError(sql.ErrConnDone)
		_, err := store.GetAllTask()
		require.Error(t, err)
	})

	t.Run("Scan Error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "description", "status", "userid"}).
			AddRow(1, "X", true, 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, status , userid FROM tasks")).
			WillReturnRows(rows)
		rows.RowError(0, errors.New("scan error"))
		_, err := store.GetAllTask()
		require.Error(t, err)
	})
}

func Test_GetTasksByUserIDTask(t *testing.T) {
	store, mock, cleanup := setup(t)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, status , userid FROM tasks where userid =?")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "description", "status", "userid"}).
				AddRow(1, "User task", true, 1))
		tasks, err := store.GetTasksByUserIDTask(1)
		require.NoError(t, err)
		require.Len(t, tasks, 1)
		require.Equal(t, 1, tasks[0].Userid)
	})

	t.Run("Query Error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, status , userid FROM tasks where userid =?")).
			WithArgs(999).
			WillReturnError(sql.ErrConnDone)
		_, err := store.GetTasksByUserIDTask(999)
		require.Error(t, err)
	})

	t.Run("Scan Error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "description", "status", "userid"}).
			AddRow(2, "B", false, 999)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, status , userid FROM tasks where userid =?")).
			WithArgs(999).
			WillReturnRows(rows)
		rows.RowError(0, errors.New("scan error"))
		_, err := store.GetTasksByUserIDTask(999)
		require.Error(t, err)
	})

}
