package comments

import (
	"gopkg.in/mgo.v2/bson"
	"reddit/pkg/user"
	"time"
)

type Comment struct {
	Created time.Time
	Author  user.User
	Body    string
	ID      bson.ObjectId
}

type CommentRepo interface {
	GetAll(bson.ObjectId) ([]*Comment, error)
	Add(*Comment, bson.ObjectId) (uint8, error)
	Delete(bson.ObjectId, bson.ObjectId)
}
