package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/SoulHardd/service-courier/internal/config"
	http2 "github.com/SoulHardd/service-courier/internal/handlers/http"
	"github.com/SoulHardd/service-courier/internal/handlers/http/courier"
	"github.com/SoulHardd/service-courier/internal/handlers/http/delivery"
	"github.com/SoulHardd/service-courier/internal/logger"
	"github.com/SoulHardd/service-courier/internal/middleware"
	"github.com/SoulHardd/service-courier/internal/rateLimiter"

	chi "github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	HttpSrv *http.Server
	r       *chi.Mux
	cfg     *config.ServerConfig
	db      *pgxpool.Pool
	lg      logger.Logger
}

func New(
	cfg *config.ServerConfig,
	pool *pgxpool.Pool,
	courier *courier.CourierController,
	delivery *delivery.DeliveryController,
	logger logger.Logger,
	rateLimiterCfg *config.RateLimiterConfig,
) *Server {
	r := chi.NewRouter()
	ipLimiter := rateLimiter.NewIPRateLimiter()
	for _, mw := range middleware.NewMiddleware(logger, ipLimiter, rateLimiterCfg) {
		r.Use(mw)
	}

	s := &Server{
		HttpSrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: r,
		},
		r:   r,
		cfg: cfg,
		db:  pool,
		lg:  logger,
	}
	s.setRoutes(courier, delivery)
	return s
}

func (s *Server) ListenAndServe() {
	go func() {
		log.Printf("Server is listening on %d", s.cfg.Port)
		if err := s.HttpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()
}

func (s *Server) setRoutes(courier *courier.CourierController, delivery *delivery.DeliveryController) {
	s.r.Get("/ping", http2.Ping)
	s.r.Head("/healthcheck", http2.HealthCheck)
	s.r.Get("/couriers", courier.GetAll)
	s.r.Route("/courier", func(r chi.Router) {
		r.Get("/{id}", courier.Get)
		r.Post("/", courier.Create)
		r.Put("/", courier.Update)
	})
	s.r.Route("/delivery", func(r chi.Router) {
		r.Post("/assign", delivery.Create)
		r.Post("/unassign", delivery.Delete)
	})
}
