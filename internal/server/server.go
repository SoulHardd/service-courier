package server

import (
	"avito/internal/config"
	"avito/internal/handlers"
	"avito/internal/handlers/courier"
	"avito/internal/handlers/delivery"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	HttpSrv *http.Server
	r       *chi.Mux
	cfg     *config.ServerConfig
	db      *pgxpool.Pool
}

func New(cfg *config.ServerConfig, pool *pgxpool.Pool, courier *courier.CourierController, delivery *delivery.DeliveryController) *Server {
	r := chi.NewRouter()
	s := &Server{
		HttpSrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: r,
		},
		r:   r,
		cfg: cfg,
		db:  pool,
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
	s.r.Get("/ping", handlers.Ping)
	s.r.Head("/healthcheck", handlers.HealthCheck)
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
