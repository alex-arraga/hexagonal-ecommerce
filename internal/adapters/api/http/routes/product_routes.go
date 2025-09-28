package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadProductRoutes(r chi.Router, h *handlers.ProductHandler) {
	r.Route("/product", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			h.ListProducts(r, w)
		})
		r.Get("/{product_id}", func(w http.ResponseWriter, r *http.Request) {
			h.FindProductById(r, w)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			h.SaveProduct(r, w)
		})
		r.Put("/{product_id}", func(w http.ResponseWriter, r *http.Request) {
			h.SaveProduct(r, w)
		})
		r.Delete("/{product_id}", func(w http.ResponseWriter, r *http.Request) {
			h.DeleteProduct(r, w)
		})
	})
}
