package rateLimit

import (
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
	ip := ipAddress.ReadUserIP(r)
	token := accessToken.ReadAccessToken(r, accessKey)
	timeExpiration, exists := os.LookupEnv("TIME_EXP")

	if !exists {
		timeExpiration = "10"
	}

	key := "limiter:ip:" + ip
	limit := 10

	if token != "" {
		key = "limiter:token:" + token
		limitValue, err := dbClient.Get(key)
		if err != nil {
			return false
		}

		limit, _ = strconv.Atoi(limitValue.(string))
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
