package user

import (
	"database/sql"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

type UserSqlRepo struct {
	DB *sql.DB
}

func NewSqlRepo(db *sql.DB) *UserSqlRepo {
	return &UserSqlRepo{DB: db}
}

var (
	ErrNoUser     = errors.New("No user found")
	ErrBadPass    = errors.New("Invalid password")
	ErrUserExists = errors.New("Username already exists")
	ErrNoMatch    = errors.New("Passwords must match")
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

func (repo *UserSqlRepo) Register(login, password, confirmation string) (*User, error) {
	if password != confirmation {
		return nil, ErrNoMatch
	}

	row := repo.DB.QueryRow("SELECT login FROM users WHERE login = ?", login)
	err := row.Scan()
	if err != sql.ErrNoRows {
		return nil, ErrUserExists
	}

	userID := bson.NewObjectId()
	_, err = repo.DB.Exec("INSERT  INTO users (`login`, `password`, `id`) VALUES (?, ?, ?)",
		login,
		password,
		userID,
	)
	if err != nil {
		return nil, err
	}

	user := &User{
		Id:       userID,
		Login:    login,
		password: password,
	}
	return user, nil
}
