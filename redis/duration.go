package redis

import (
	"time"
)

const (
	SESSION_TTL     time.Duration = 30 * time.Minute
	AUTH_TOKEN_TTL  time.Duration = 15 * time.Minute
	CACHE_TTL_SHORT time.Duration = 1 * time.Minute
	CACHE_TTL_LONG  time.Duration = 1 * time.Hour
)
