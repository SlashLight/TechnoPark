package handlers

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"reddit/pkg/items"
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

func (h *ItemsHandler) ListByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorID := vars["user_id"]
	elems, err := h.ItemsRepo.GetByAuthor(authorID)

}
