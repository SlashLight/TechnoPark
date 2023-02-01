package comments

import (
	"gopkg.in/mgo.v2/bson"
	"reddit/pkg/items"
)

func (repo *items.ItemMongoRepository) GetAll(postId bson.ObjectId) ([]*Comment, error) {
	post := &Item{}
	err := repo.
}
