package api

import (
	"sync"
	"time"

	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

type RustMapsClientBase interface {
	GetStatus(log *zap.Logger, m *types.Map) (*RustMapsStatusResponse, error)
	SetApiKey(apiKey string)
	GetLimits(log *zap.Logger) (*RustMapsLimitsResponse, error)
	GenerateCustom(log *zap.Logger, m *types.Map) (*RustMapsGenerateResponse, error)
	GenerateProcedural(log *zap.Logger, m *types.Map) (*RustMapsGenerateResponse, error)
}

type RustMapsClient struct {
	RustMapsClientBase
	ApiUrl      string
	apiKey      string
	rateLimiter *RateLimiter
}

func NewRustMapsClient(apiKey string) RustMapsClientBase {
	return &RustMapsClient{
		ApiUrl:      "https://api.rustmaps.com/v4",
		apiKey:      apiKey,
		rateLimiter: NewRateLimiter(60),
	}
}

func (r *RustMapsClient) SetApiKey(apiKey string) {
	r.apiKey = apiKey
}

// RateLimiter manages API request timing
type RateLimiter struct {
	callsPerMinute int
	interval       time.Duration
	lastCall       time.Time
	mu             sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(callsPerMinute int) *RateLimiter {
	return &RateLimiter{
		callsPerMinute: callsPerMinute,
		interval:       time.Minute / time.Duration(callsPerMinute),
	}
}

// Wait ensures enough time has passed since the last call
func (r *RateLimiter) Wait() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if !r.lastCall.IsZero() {
		timePassed := now.Sub(r.lastCall)
		if timePassed < r.interval {
			time.Sleep(r.interval - timePassed)
		}
	}
	r.lastCall = time.Now()
}
