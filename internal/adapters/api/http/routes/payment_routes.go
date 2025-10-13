package routes

import (
	"go-ecommerce/internal/adapters/api/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadPaymentRoutes(r chi.Router, h *handlers.PaymentHandler) {
	r.Route("/payment", func(r chi.Router) {
		r.Post("/mp", func(w http.ResponseWriter, r *http.Request) {
			h.StartTransaction(r, w)
		})
		r.Post("/mp/webhook", func(w http.ResponseWriter, r *http.Request) {
			h.NotificationWebhook(r, w)
		})
	})
}
