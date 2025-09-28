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
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			h.ListCategories(r, w)
		})
		r.Get("/{category_id}", func(w http.ResponseWriter, r *http.Request) {
			h.FindCategoryById(r, w)
		})
		r.Delete("/{category_id}", func(w http.ResponseWriter, r *http.Request) {
			h.DeleteCategory(r, w)
		})
	})
}
