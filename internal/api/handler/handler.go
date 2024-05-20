package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	loyaltyMiddleware "e1m0re/loyalty-srv/internal/api/handler/middleware"
	"e1m0re/loyalty-srv/internal/service"
)

type Handler interface {
	AddOrder(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
	GetOrders(w http.ResponseWriter, r *http.Request)
	GetWithdrawals(w http.ResponseWriter, r *http.Request)
	NewRouter() *chi.Mux
	SignIn(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	WritingOff(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	UserService    service.UsersService
	OrderService   service.OrdersService
	AccountService service.AccountsService
}

func NewHandler(userService service.UsersService, orderService service.OrdersService, accountService service.AccountsService) Handler {
	return &handler{
		UserService:    userService,
		OrderService:   orderService,
		AccountService: accountService,
	}

}

func (handler handler) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Compress(5, "text/html", "application/json"))

	// Public functions
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handler.SignUp)
		r.Post("/api/user/login", handler.SignIn)
	})

	// Private functions
	r.Group(func(r chi.Router) {
		r.Use(loyaltyMiddleware.PrivateRoute)
		r.Route("/", func(r chi.Router) {
			r.Route("/api/user", func(r chi.Router) {
				r.Route("/orders", func(r chi.Router) {
					r.Get("/", handler.GetOrders)
					r.Post("/", handler.AddOrder)
				})
				r.Route("/balance", func(r chi.Router) {
					r.Get("/", handler.GetBalance)
					r.Post("/withdraw", handler.WritingOff)
				})
				r.Get("/withdrawals", handler.GetWithdrawals)
			})
		})
	})

	return r
}
