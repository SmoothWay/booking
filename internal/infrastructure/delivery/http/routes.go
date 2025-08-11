package http

import (
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/handler"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router *chi.Mux, handler *handler.Handler) {

	router.Route("/api/v1/booking", func(r chi.Router) {
		r.Get("/", handler.GetBookings)
		// r.Post("/", handler.CreateBooking)
		// r.Put("/{id}", handler.UpdateBooking)
		// r.Delete("/{id}", handler.DeleteBooking)
	})

	router.Route("/api/v1/user", func(r chi.Router) {
		r.Get("/", handler.GetUsers)
	})

	router.Route("/api/v1/unit", func(r chi.Router) {
		r.Get("/", handler.GetUnits)
	})

	router.Route("/probes", func(r chi.Router) {
		r.Get("/health", handler.Health)
		r.Get("/liveness", handler.Liveness)
		r.Get("/readiness", handler.Readiness)
	})
}
