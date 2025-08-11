package http

import (
	"fmt"
	"net/http"

	"github.com/SmoothWay/booking/internal/config"
	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	cfg    config.HTTP
}

func NewServer(cfg config.HTTP, bookingUsecase domain.BookingUsecase, userUsecase domain.UserUsecase, unitUsecase domain.UnitUsecase) *Server {
	router := chi.NewRouter()
	handler := handler.NewHandler(bookingUsecase, userUsecase, unitUsecase)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)

	SetupRoutes(router, handler)

	return &Server{router: router, cfg: cfg}
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), s.router)
}
