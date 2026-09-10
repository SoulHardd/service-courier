package middleware

import (
	"net/http"

	"github.com/SoulHardd/service-courier/internal/config"
	"github.com/SoulHardd/service-courier/internal/logger"
	"github.com/SoulHardd/service-courier/internal/rateLimiter"
)

func NewMiddleware(log logger.Logger, ipLimiter *rateLimiter.IPRateLimiter, rateLimiterCfg *config.RateLimiterConfig) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		RateLimitMiddleware(ipLimiter, log, rateLimiterCfg),
		MetricsAndLogsMiddleware(log),
	}
}
