package user

import (
	"Task_Manager/model/user"
	"database/sql"
)

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}

func (us *UserStore) CreateUser(user user.User) (user.User, error) {
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	result, err := us.DB.Exec(query, user.Name, user.Email)

	if err != nil {
		return user, err
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)

	return user, nil
}

func (us *UserStore) GetByIDUser(id int) (user.User, error) {
	var user user.User

	query := "SELECT id, name, email FROM users WHERE id = ?"
	err := us.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return user, err
	}

	return user, err
}

func (us *UserStore) DeleteUser(id int) error {
	_, err := us.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (us *UserStore) GetAllUser() ([]user.User, error) {
	query := "SELECT id, name, email FROM users"
	rows, err := us.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var users []user.User

	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil

}
