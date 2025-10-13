package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/core/ports"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	srv ports.PaymentProvider
}

func NewPaymentHandler(paymentService ports.PaymentProvider) *PaymentHandler {
	return &PaymentHandler{srv: paymentService}
}

func (ph *PaymentHandler) StartTransaction(r *http.Request, w http.ResponseWriter) {
	// Verify HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get order id and converts in uuid
	orderID := chi.URLParam(r, "order_id")
	if orderID == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "OrderID is required to starting transaction")
		return
	}

	uuid, err := uuid.Parse(orderID)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Order id must be a valid uuid: %s", err))
		return
	}

	// If operations it's ok, return a redirect_url to can pay
	redirectUrl, err := ph.srv.CreatePayment(r.Context(), uuid)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error processing payment request: %s", err))
		return
	}

	// Return tre redirect_url that mercado pago provided
	httpdtos.RespondJSON(w, http.StatusOK, "Successfully payment request processed", redirectUrl)
}

func (ph *PaymentHandler) NotificationWebhook(r *http.Request, w http.ResponseWriter) {
	// Verify HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract id and topic from mercado pago request
	topic := r.URL.Query().Get("topic")
	id := r.URL.Query().Get("id")

	if id == "" && topic == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "Parameters id or topic not found")
		return
	}

	err := ph.srv.VerifyPayment(r.Context(), &id, &topic)
	if err != nil {
		httpdtos.RespondJSON(w, http.StatusOK, "Error updating order with payment data in webhook", err)
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Order successfully updated with payment data", nil)
}
