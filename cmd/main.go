package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"rateLimiter/internal/middleware"
	"rateLimiter/internal/repository/redisRepository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisClient := redisRepository.Config()
	redisRepo := redisRepository.NewRedisRepository(redisClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	wrappedMux := middleware.RateLimitMiddleware(redisRepo)(mux)

	http.ListenAndServe(":8080", wrappedMux)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
