package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"rateLimiter/internal/middleware"
	"rateLimiter/internal/repository/redisRepository"
)

func main() {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	redisRepo := redisRepository.NewRedisRepository(redisClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/", yourMainHandler)

	wrappedMux := middleware.RateLimitMiddleware(redisRepo)(mux)

	http.ListenAndServe(":8080", wrappedMux)
}

func yourMainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
