package session

import (
	"database/sql"
	"net/http"
	"time"
)

type SessionManager struct {
	DB *sql.DB
}

func NewSessionManager(db *sql.DB) *SessionManager {
	return &SessionManager{DB: db}
}

func (sm *SessionManager) Check(r *http.Request) (*Session, error) {
	sess := &Session{}
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, ErrNoAuth
	}

	row := sm.DB.QueryRow("SELECT id, userid FROM sessions WHERE id = ?", sessionCookie.Value)
	err = row.Scan(&sess.ID, &sess.UserID)
	if err == sql.ErrNoRows {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func (sm *SessionManager) Create(w http.ResponseWriter, userID uint32) (*Session, error) {
	sess := NewSession(userID)

	_, err := sm.DB.Exec("INSERT INTO sessions ('id', 'userid') VALUES (?, ?)", sess.ID, sess.UserID)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Path:    "/",
		Expires: time.Now().Add(90 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)
	return sess, nil
}

func (sm *SessionManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	_, err = sm.DB.Exec("DELETE FROM sessions WHERE id = ?", sess.ID)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}
