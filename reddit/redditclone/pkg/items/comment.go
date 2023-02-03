package items

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Comment struct {
	Created time.Time
	Author  Author
	Body    string
	ID      bson.ObjectId
}

type CommentRepo interface {
	GetAll(bson.ObjectId) ([]*Comment, error)
	Add(*Comment, bson.ObjectId) (uint8, error)
	Delete(bson.ObjectId, bson.ObjectId)
}
