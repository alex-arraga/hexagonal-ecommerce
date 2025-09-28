package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadUserRoutes(r chi.Router, h *handlers.UserHandler) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/find/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			h.FindUserById(r, w)
		})
		r.Get("/find", func(w http.ResponseWriter, r *http.Request) {
			h.FindUserByEmail(r, w)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			h.ListUsers(r, w)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			h.SaveUser(r, w)
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			h.SaveUser(r, w)
		})
		r.Delete("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			h.DeleteUser(r, w)
		})
	})
}
