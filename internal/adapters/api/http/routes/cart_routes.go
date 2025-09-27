package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadCartRoutes(r chi.Router, h *handlers.CartHandler) {
	r.Route("/user/{user_id}/cart", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			h.GetProductsFromCart(r, w)
		})
		r.Post("/{product_id}", func(w http.ResponseWriter, r *http.Request) {
			h.AddProductToCart(r, w)
		})
		r.Put("/{product_id}", func(w http.ResponseWriter, r *http.Request) {
			h.AddProductToCart(r, w)
		})
		r.Delete("/{product_id}", func(w http.ResponseWriter, r *http.Request) {
			h.RemoveItemFromCart(r, w)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			h.ClearCart(r, w)
		})
	})
}
