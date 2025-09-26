package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadOrderRoutes(r chi.Router, h *handlers.OrderHandler) {
	r.Route("/order", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			h.SaveOrder(r, w)
		})
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.SaveOrder(r, w)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.GetOrderByID(r, w)
		})
	})
}
