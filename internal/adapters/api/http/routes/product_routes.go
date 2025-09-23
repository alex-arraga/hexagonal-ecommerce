package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadProductRoutes(r chi.Router, h *handlers.ProductHandler) {
	r.Route("/product", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			h.SaveProduct(r, w)
		})
	})
}
