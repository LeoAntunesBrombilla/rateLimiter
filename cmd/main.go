package main

import (
	"github.com/go-redis/redis/v8"
	"net/http"
	"os"
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
	addr := os.Getenv("REDIS_ADDR")

	if addr == "" {
		addr = "localhost:6379"
	}

	redisAddr, exists := os.LookupEnv("REDIS_ADDR")

	if !exists {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	limiter := NewRateLimiter(client, 10)

	http.Handle("/", limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})))

	http.ListenAndServe(":8080", nil)
}
