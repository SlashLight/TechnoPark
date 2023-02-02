package handlers

import (
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"reddit/pkg/session"
	"reddit/pkg/user"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo user.UserRepo
	Sessions *session.SessionManager
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Authorize(r.FormValue("username"), r.FormValue("password"))
	if err == user.ErrNoUser {
		http.Error(w, `no user`, http.StatusBadRequest)
		return
	}
	if err == user.ErrBadPass {
		http.Error(w, `bad pass`, http.StatusBadRequest)
		return
	}

	sess, _ := h.Sessions.Create(w, u.Id)
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", 302)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Register(r.FormValue("username"), r.FormValue("password"), r.FormValue("confirm password"))
	if err == user.ErrNoMatch {
		http.Error(w, `passwords must match`, http.StatusBadRequest)
		return
	}
	if err == user.ErrUserExists {
		http.Error(w, `username already exists`, http.StatusBadRequest)
		return
	}

	sess, _ := h.Sessions.Create(w, u.Id)
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", 302)
}
