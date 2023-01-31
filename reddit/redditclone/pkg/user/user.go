package user

type User struct {
	Id       uint32
	Login    string
	password string
}

type UserRepo interface {
	Authorize(login, password string) (*User, error)
}
