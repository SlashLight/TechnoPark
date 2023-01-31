package user

import (
	"database/sql"
	"errors"
)

type UserSqlRepo struct {
	DB *sql.DB
}

func NewSqlRepo(db *sql.DB) *UserSqlRepo {
	return &UserSqlRepo{DB: db}
}

var (
	ErrNoUser  = errors.New("No user found")
	ErrBadPass = errors.New("Invalid password")
)

func (repo *UserSqlRepo) Authorize(login, password string) (*User, error) {
	user := &User{}
	row := repo.DB.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login)
	err := row.Scan(&user.Id, &user.Login, &user.password)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}

	if user.password != password {
		return nil, ErrBadPass
	}
	
	return user, nil
}
