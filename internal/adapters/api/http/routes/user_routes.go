package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadUserRoutes(r chi.Router, h *handlers.UserHandler) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			h.SaveUser(w, r)
		})
	})
}
