package middleware

import (
	"net/http"
	"reddit/pkg/session"
)

var (
	noAuthUrls = map[string]struct{}{
		"/login": struct{}{},
	}

	noSessUrls = map[string]struct{}{
		"/": struct{}{},
	}
)

func Auth(sm *session.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		_, canBeWithoutSess := noSessUrls[r.URL.Path]
		if err != nil && !canBeWithoutSess {
			http.Redirect(w, r, "/", 302)
			return
		}

		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
