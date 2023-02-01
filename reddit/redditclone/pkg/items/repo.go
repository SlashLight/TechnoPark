package items

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	validCategories = map[string]struct{}{
		"music":       struct{}{},
		"funny":       struct{}{},
		"videos":      struct{}{},
		"programming": struct{}{},
		"news":        struct{}{},
		"fashion":     struct{}{},
	}
)

type ItemMongoRepository struct {
	Sess  *mgo.Session
	Items *mgo.Collection
}

func NewMongoRepo(sess *mgo.Session, coll *mgo.Collection) *ItemMongoRepository {
	return &ItemMongoRepository{Sess: sess, Items: coll}
}

func (repo *ItemMongoRepository) GetAll() ([]*Item, error) {
	items := []*Item{}

	err := repo.Items.Find(bson.M{}).All(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ItemMongoRepository) GetByID(id bson.ObjectId) (*Item, error) {
	post := &Item{}
	err := repo.Items.Find(bson.M{"id": id}).One(&post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

var ErrWrongCategory = errors.New("Invalid category")

func (repo *ItemMongoRepository) GetByCategory(categ string) ([]*Item, error) {
	if _, ok := validCategories[categ]; !ok {
		return nil, ErrWrongCategory
	}

	items := []*Item{}
	err := repo.Items.Find(bson.M{"category": categ}).All(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ItemMongoRepository) GetByAuthor(author string) ([]*Item, error) {
	items := []*Item{}
	err := repo.Items.Find(bson.M{"author": author}).All(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ItemMongoRepository) Upvote(id bson.ObjectId) (uint16, error) {
	post := &Item{}
	err := repo.Items.Find(bson.M{"id": id}).One(&post)
	if err != nil {
		return 0, err
	}

	newVote := &Vote{User: post.Author.Id, Vote: 1}
	post.Votes = append(post.Votes, newVote)
	post.Score = post.Score + 1
	post.UpvotePercentage = post.Score / uint16(len(post.Votes))

	err = repo.Items.Update(bson.M{"id": id}, &post)
	if err == mgo.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return post.Score, nil
}

func (repo *ItemMongoRepository) Downvote(id bson.ObjectId) (uint16, error) {
	post := &Item{}
	err := repo.Items.Find(bson.M{"id": id}).One(&post)
	if err != nil {
		return 0, err
	}

	newVote := &Vote{User: post.Author.Id, Vote: -1}
	post.Votes = append(post.Votes, newVote)
	post.Score = post.Score - 1
	post.UpvotePercentage = post.Score / uint16(len(post.Votes))

	err = repo.Items.Update(bson.M{"id": id}, &post)
	if err == mgo.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return post.Score, nil
}

func (repo *ItemMongoRepository) Add(newItem *Item) (uint8, error) {
	err := repo.Items.Insert(newItem)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (repo *ItemMongoRepository) Delete(id bson.ObjectId) (uint8, error) {
	err := repo.Items.Remove(bson.M{"id": id})
	if err == mgo.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return 1, nil
}

func (repo *ItemMongoRepository) AddComment(postId bson.ObjectId, comment *Comment) (uint8, error) {
	post := &Item{}
	err := repo.Items.Find(bson.M{"id": postId}).One(&post)
	if err != nil {
		return 0, err
	}

	post.Comments = append(post.Comments, comment)
	err = repo.Items.Update(bson.M{"id": postId}, &post)
	if err == mgo.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return 1, nil
}

func (repo *ItemMongoRepository) DeleteComment(postId bson.ObjectId, commentId bson.ObjectId) (uint8, error) {
	post := &Item{}
	err := repo.Items.Find(bson.M{"id": postId}).One(&post)
	if err != nil {
		return 0, err
	}

	ind := -1
	for idx, comment := range post.Comments {
		if comment.ID == commentId {
			ind = idx
			break
		}
	}

	if ind == -1 {
		ErrWrongCom := errors.New("Invalid comment ID")
		return 0, ErrWrongCom
	}
	post.Comments = append(post.Comments[:ind], post.Comments[ind+1:]...)

	err = repo.Items.Update(bson.M{"id": postId}, &post)
	if err == mgo.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return 1, nil
}
