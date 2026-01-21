package middleware

import (
	"avito/internal/config"
	"avito/internal/logger"
	"avito/internal/metrics"
	"avito/internal/rateLimiter"
	"net/http"
)

func RateLimitMiddleware(ipLimiter *rateLimiter.IPRateLimiter, log logger.Logger, cfg *config.RateLimiterConfig) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip := r.RemoteAddr
			limiter := ipLimiter.Get(ip, cfg.MaxTokens, cfg.RefillRate)

			if !limiter.Allow() {
				metrics.RateLimitExceededTotal.Inc()

				log.Warn("rate limit exceeded",
					logger.Field{Key: "ip", Value: ip},
					logger.Field{Key: "path", Value: r.URL.Path},
				)

				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
