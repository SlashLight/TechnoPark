package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"reddit/pkg/items"
	"reddit/pkg/session"
	"time"
)

type ItemsHandler struct {
	Tmpl      *template.Template
	ItemsRepo items.ItemRepo
	Logger    *zap.SugaredLogger
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	elems, err := h.ItemsRepo.GetAll()
	if err != nil {
		h.Logger.Error("GetAll err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(elems)
	w.Write(respJSON)
}

func (h *ItemsHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	elems, err := h.ItemsRepo.GetByCategory(category)
	if err == items.ErrWrongCategory {
		h.Logger.Error("Wrong category", err)
		http.Error(w, `Wrong category`, http.StatusBadRequest)
		return
	} else if err != nil {
		h.Logger.Error("GetByCategory err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(elems)
	w.Write(respJSON)
}

/*func (h *ItemsHandler) ListByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorID := vars["user_id"]
	elems, err := h.ItemsRepo.GetByAuthor(authorID)

} */

func (h *ItemsHandler) Add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	item := new(items.Item)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(item, r.PostForm)
	if err != nil {
		h.Logger.Error("Form err", err)
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		h.Logger.Error("Sess err", err)
		http.Error(w, `Auth err`, http.StatusBadRequest)
		return
	}
	author := &items.Author{
		Username: sess.Username,
		ID:       bson.NewObjectId(),
	}

	item.Author = author
	item.Created = time.Now()
	item.ID = bson.NewObjectId()

	inserted, err := h.ItemsRepo.Add(item)
	if err != nil {
		h.Logger.Error("Db err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert post: %v", inserted)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(item)
	w.Write(respJSON)
}

func (h *ItemsHandler) PostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ok := bson.IsObjectIdHex(vars["id"])
	if !ok {
		http.Error(w, `{"error": "bad id"`, http.StatusBadGateway)
		return
	}

	id := bson.ObjectId(vars["id"])
	item, err := h.ItemsRepo.GetByID(id)
	if err != nil {
		h.Logger.Error("GetByID err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(item)
	w.Write(respJSON)
}

func (h *ItemsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `Auth err`, http.StatusBadRequest)
		return
	}

	comment := new(items.Comment)
	comment.Body = r.FormValue("comment")
	comment.Created = time.Now()
	comment.Author = items.Author{
		Username: sess.Username,
		ID:       sess.UserID,
	}
	comment.ID = bson.NewObjectId()

	vars := mux.Vars(r)
	ok := bson.IsObjectIdHex(vars["post_id"])
	if !ok {
		http.Error(w, `{"error": "bad post id"`, http.StatusBadGateway)
		return
	}
	postID := bson.ObjectId(vars["post_id"])
	inserted, err := h.ItemsRepo.AddComment(postID, comment)
	if err != nil {
		h.Logger.Error("AddComment err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("insert comment: %v", inserted)

	newPost, _ := h.ItemsRepo.GetByID(postID)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(newPost)
	w.Write(respJSON)
}

func (h *ItemsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ok := bson.IsObjectIdHex(vars["post_id"])
	if !ok {
		http.Error(w, `{"error": "bad post id"`, http.StatusBadGateway)
		return
	}
	ok = bson.IsObjectIdHex(vars["comment_id"])
	if !ok {
		http.Error(w, `{"error": "bad comment id"`, http.StatusBadGateway)
		return
	}

	postID, commentID := bson.ObjectId(vars["post_id"]), bson.ObjectId(vars["comment_id"])
	ok, err := h.ItemsRepo.DeleteComment(postID, commentID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, `{"error": "bad post id"}`, http.StatusInternalServerError)
		return
	}

	newPost, _ := h.ItemsRepo.GetByID(postID)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(newPost)
	w.Write(respJSON)
}

func (h *ItemsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ok := bson.IsObjectIdHex(vars["post_id"])
	if !ok {
		http.Error(w, `{"error": "bad post id"`, http.StatusBadGateway)
		return
	}

	postID := bson.ObjectId(vars["post_id"])
	ok, err := h.ItemsRepo.Delete(postID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, `{"error": "bad post id"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]string{
		"message": "success",
	})
	w.Write(respJSON)
}

func (h *ItemsHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ok := bson.IsObjectIdHex(vars["post_id"])
	if !ok {
		http.Error(w, `{"error": "bad post id"`, http.StatusBadRequest)
		return
	}

	postID := bson.ObjectId(vars["post_id"])
	score, err := h.ItemsRepo.Upvote(postID)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("new post rating: %v", score)

	newPost, _ := h.ItemsRepo.GetByID(postID)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(newPost)
	w.Write(respJSON)
}

func (h *ItemsHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ok := bson.IsObjectIdHex(vars["post_id"])
	if !ok {
		http.Error(w, `{"error": "bad post id"`, http.StatusBadRequest)
		return
	}

	postID := bson.ObjectId(vars["post_id"])
	score, err := h.ItemsRepo.Downvote(postID)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("new post rating: %v", score)

	newPost, _ := h.ItemsRepo.GetByID(postID)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(newPost)
	w.Write(respJSON)
}
