package rateLimit

import (
	"log"
	"net/http"
	"os"
	"rateLimiter/internal/repository"
	"rateLimiter/pkg/headerInfo/accessToken"
	"rateLimiter/pkg/headerInfo/ipAddress"
	"strconv"
	"time"
)

var requestTimestamps = make(map[string]int64)
var requestCounts = make(map[string]int)

func RateLimit(r *http.Request, accessKey string, dbClient repository.Database) bool {
	limit := 10
	ip := ipAddress.ReadUserIP(r)
	key := "limiter:ip:" + ip
	apiKeyName, reqLimit := accessToken.ReadAccessToken(r, accessKey)
	timeExpiration, has := os.LookupEnv("TIME_EXP")
	if !has {
		timeExpiration = "10"
	}

	err := dbClient.Set(apiKeyName, reqLimit)
	if err != nil {
		log.Printf("Fail to save on DB: %v", err)
	}

	if reqLimit != "" {
		key = apiKeyName
		limitValue, err := dbClient.Get(key)
		if err != nil {
			log.Printf("Error getting rate limit for key %s: %v", key, err)
			return false
		}

		if limitStr, ok := limitValue.(string); ok {
			limit, _ = strconv.Atoi(limitStr)
		} else {
			log.Printf("Invalid rate limit format for key %s", key)
			return false
		}
	} else {
		key = "limiter:ip:" + ip
		limit = 10
	}

	currentTime := time.Now().Unix()
	windowDuration, err := strconv.ParseInt(timeExpiration, 10, 64)

	if err != nil {
		panic("Fatal error during conversion of token")
	}

	if start, exists := requestTimestamps[key]; exists {
		if currentTime-start < windowDuration {
			requestCounts[key]++
			if requestCounts[key] > limit {
				return false
			}
		} else {
			requestTimestamps[key] = currentTime
			requestCounts[key] = 1
		}
	} else {
		requestTimestamps[key] = currentTime
		requestCounts[key] = 1
	}

	return true
}
