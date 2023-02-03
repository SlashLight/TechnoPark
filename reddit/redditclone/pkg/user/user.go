package user

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id       bson.ObjectId
	Login    string
	password string
}

type UserRepo interface {
	Authorize(login, password string) (*User, error)
	Register(login, password, confirmation string) (*User, error)
}
