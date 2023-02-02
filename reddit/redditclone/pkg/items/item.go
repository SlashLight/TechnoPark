package items

import (
	"reddit/pkg/user"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Vote struct {
	User uint32
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
	Title            string        `bson:"title"`
	Url              string        `bson:"url"`
	Author           user.User     `bson:"author"`
	Category         string        `bson:"category"`
	Text             string        `bson:"text"`
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
	Add(*Item) (uint8, error)
	Delete(bson.ObjectId) (uint8, error)

	AddComment(bson.ObjectId, *Comment) (uint8, error)
	DeleteComment(bson.ObjectId, bson.ObjectId) (uint8, error)
}
