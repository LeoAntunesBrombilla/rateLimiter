package middleware

import (
	"github.com/LeoAntunesBrombilla/rateLimiter/internal/repository"
	"github.com/LeoAntunesBrombilla/rateLimiter/pkg/rateLimit"
	"net/http"
)

func RateLimitMiddleware(dbClient repository.Database) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rateLimit.RateLimit(r, "ACCESS_TOKEN", dbClient) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
