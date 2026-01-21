package middleware

import (
	"avito/internal/config"
	"avito/internal/logger"
	"avito/internal/rateLimiter"
	"net/http"
)

func NewMiddleware(log logger.Logger, ipLimiter *rateLimiter.IPRateLimiter, rateLimiterCfg *config.RateLimiterConfig) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		RateLimitMiddleware(ipLimiter, log, rateLimiterCfg),
		MetricsAndLogsMiddleware(log),
	}
}
