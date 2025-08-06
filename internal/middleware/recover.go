package middlewares

import (
	"log/slog"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.Error("Recovered from panic", slog.Any("error", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
