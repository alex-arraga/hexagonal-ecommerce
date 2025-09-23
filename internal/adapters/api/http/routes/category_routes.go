package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadCategoryRoutes(r chi.Router, h *handlers.CategoryHandler) {
	r.Route("/category", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			h.SaveCategory(r, w)
		})
	})
}
