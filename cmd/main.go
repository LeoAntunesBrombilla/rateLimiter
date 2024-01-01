package main

import (
	"github.com/go-redis/redis/v8"
	"net/http"
	"rateLimiter/internal/repository/redisRepository"
)

type RateLimiter struct {
	redisClient *redis.Client
	rateLimit   int
}

func NewRateLimiter(client *redis.Client, rateLimit int) *RateLimiter {
	return &RateLimiter{
		redisClient: client,
		rateLimit:   rateLimit,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func main() {
	client := redisRepository.Config()

	limiter := NewRateLimiter(client, 10)

	http.Handle("/", limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})))

	http.ListenAndServe(":8080", nil)
}
