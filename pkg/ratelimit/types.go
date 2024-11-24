package ratelimit

import "context"

//go:generate mockgen -source=./types.go -package=limitmocks -destination=mocks/limiter.mock.go Limiter
type Limiter interface {
	// Limit 要不要限流
	Limit(ctx context.Context, key string) (bool, error)
}
