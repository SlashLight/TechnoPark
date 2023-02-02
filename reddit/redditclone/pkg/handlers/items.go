package handlers

import (
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

	//EXECUTE TEMPLATE
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

	//EXECUTE TEMPLATE
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
	h.Logger.Infof("Inserted %v items", inserted)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ItemsHandler) PostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ok := bson.IsObjectIdHex(vars["id"])
	if !ok {
		http.Error(w, `{"error": "bad id"`, http.StatusBadGateway)
		return
	}

	id := vars["id"]
	item, err := h.ItemsRepo.GetByID(id)
	if err != nil {
		h.Logger.Error("GetByID err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	//TEMPLATE
}

func (h *ItemsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	
}
