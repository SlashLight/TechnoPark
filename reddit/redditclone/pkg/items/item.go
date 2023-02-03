package items

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Vote struct {
	User bson.ObjectId
	Vote int8
}

type Author struct {
	Username string
	ID       bson.ObjectId
}

type Item struct {
	Score            uint16        `bson:"score"`
	Views            uint16        `bson:"views"`
	Type             string        `bson:"type"`
	Title            string        `bson:"title" schema:"title"`
	Url              string        `bson:"url, omitempty" schema:"url"`
	Author           *Author       `bson:"author"`
	Category         string        `bson:"category" schema:"category"`
	Text             string        `bson:"text, omitempty" schema:"text"`
	Votes            []*Vote       `bson:"votes"`
	Comments         []*Comment    `bson:"comments"`
	Created          time.Time     `bson:"created"`
	UpvotePercentage uint16        `bson:"upvotePercentage"`
	ID               bson.ObjectId `bson:"ID"`
}

type ItemRepo interface {
	GetAll() ([]*Item, error)
	GetByID(bson.ObjectId) (*Item, error)
	GetByCategory(string) ([]*Item, error)
	GetByAuthor(string) ([]*Item, error)
	Upvote(bson.ObjectId) (uint16, error)
	Downvote(bson.ObjectId) (uint16, error)
	Add(*Item) (bool, error)
	Delete(bson.ObjectId) (bool, error)

	AddComment(bson.ObjectId, *Comment) (bool, error)
	DeleteComment(postId bson.ObjectId, commentId bson.ObjectId) (bool, error)
}
