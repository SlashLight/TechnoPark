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
	Score            uint16        `json:"score" bson:"score"`
	Views            uint16        `json:"views" bson:"views"`
	Type             string        `json:"type" bson:"type"`
	Title            string        `json:"title" schema:"title" bson:"title"`
	Url              string        `json:"url, omitempty" schema:"url" bson:"url, omitempty"`
	Author           *Author       `json:"author" bson:"author"`
	Category         string        `json:"category" schema:"category" bson:"category"`
	Text             string        `json:"text, omitempty" schema:"text" bson:"text, omitempty"`
	Votes            []*Vote       `json:"votes" bson:"votes"`
	Comments         []*Comment    `json:"comments" bson:"comments"`
	Created          time.Time     `json:"created" bson:"created"`
	UpvotePercentage uint16        `json:"upvotePercentage" bson:"upvotePercentage"`
	ID               bson.ObjectId `json:"id" bson:"id"`
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
