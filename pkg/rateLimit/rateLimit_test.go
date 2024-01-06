package rateLimit

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Get(key string) (interface{}, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

func (m *MockRedisRepository) Set(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

// TestRateLimit verifica se a função RateLimit permite uma requisição
// quando o limite de taxa não é excedido. A função mock 'Get' é configurada
// para retornar um limite alto, garantindo que a requisição seja permitida.
func TestRateLimit(t *testing.T) {
	mockRepo := new(MockRedisRepository)

	mockRepo.On("Get", "limiter:token:example_token").Return("100", nil)
	r, _ := http.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:1234"
	r.Header.Set("ACCESS_TOKEN", "example_token")

	result := RateLimit(r, "ACCESS_TOKEN", mockRepo)

	assert.True(t, result, "Expected RateLimit to return true")

	mockRepo.AssertExpectations(t)
}

// TestRateLimitExceeded verifica se a função RateLimit bloqueia uma requisição
// adicional quando o limite de taxa é excedido. A função mock 'Get' é configurada
// para retornar um limite de 10 requisições, e o loop simula essas 10 requisições.
// A requisição subsequente deve ser bloqueada, indicando que o limite foi excedido.
func TestRateLimitExceeded(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	mockRepo.On("Get", "limiter:token:example_token").Return("10", nil)

	r, _ := http.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:1234"
	r.Header.Set("ACCESS_TOKEN", "example_token")

	for i := 0; i < 10; i++ {
		_ = RateLimit(r, "ACCESS_TOKEN", mockRepo)
	}

	result := RateLimit(r, "ACCESS_TOKEN", mockRepo)

	assert.False(t, result, "Expected RateLimit to return false, indicating the rate limit is exceeded")

	mockRepo.AssertExpectations(t)
}

// TestRateLimitWithIP verifica se a função RateLimit permite uma requisição
// com base no endereço IP quando o limite de taxa não é excedido.
func TestRateLimitWithIP(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	ip := "127.0.0.1:1234"

	r, _ := http.NewRequest("GET", "/", nil)
	r.RemoteAddr = ip

	requestCounts = make(map[string]int)
	requestTimestamps = make(map[string]int64)

	result := RateLimit(r, "", mockRepo)
	assert.True(t, result, "Expected RateLimit to return true for IP-based rate limiting")

	mockRepo.AssertNotCalled(t, "Get", "limiter:ip:"+ip)
}

// TestRateLimitExceededWithIP verifica se a função RateLimit bloqueia uma requisição
// adicional com base no endereço IP quando o limite de taxa é excedido.
func TestRateLimitExceededWithIP(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	ip := "127.0.0.1:1234"

	r, _ := http.NewRequest("GET", "/", nil)
	r.RemoteAddr = ip

	requestCounts = make(map[string]int)
	requestTimestamps = make(map[string]int64)

	for i := 0; i < 10; i++ {
		_ = RateLimit(r, "", mockRepo)
	}

	result := RateLimit(r, "", mockRepo)
	assert.False(t, result, "Expected RateLimit to return false for IP-based rate limiting, indicating the rate limit is exceeded")

	mockRepo.AssertNotCalled(t, "Get", "limiter:ip:"+ip)
}
