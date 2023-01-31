package session

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
)

type Session struct {
	ID     string
	UserID uint32
}

func NewSession(userId uint32) *Session {
	randId := make([]byte, 16)
	rand.Read(randId)

	return &Session{
		ID:     fmt.Sprintf("%x", randId),
		UserID: userId,
	}
}

var ErrNoAuth = errors.New("No session found")

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
