package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	loyaltyMiddleware "e1m0re/loyalty-srv/internal/api/handler/middleware"
	"e1m0re/loyalty-srv/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}

}

func (h *Handler) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Compress(5, "text/html", "application/json"))

	// Public functions
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", h.SignUp)
		r.Post("/api/user/login", h.SignIn)
	})

	// Private functions
	r.Group(func(r chi.Router) {
		r.Use(loyaltyMiddleware.PrivateRoute)
		r.Route("/", func(r chi.Router) {
			r.Route("/api/user", func(r chi.Router) {
				r.Route("/orders", func(r chi.Router) {
					r.Get("/", h.GetOrders)
					r.Post("/", h.AddOrder)
				})
				r.Route("/balance", func(r chi.Router) {
					r.Get("/", h.GetBalance)
					r.Post("/withdraw", h.WritingOff)
				})
				r.Get("/withdrawals", h.GetWithdrawals)
			})
		})
	})

	return r
}
