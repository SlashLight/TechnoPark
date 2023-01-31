package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func AccessLog(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Infow("New request",
			"Method", r.Method,
			"Remote addr", r.RemoteAddr,
			"Url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}
