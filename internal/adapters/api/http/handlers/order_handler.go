package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/adapters/api/http/utils"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	srv ports.OrderService
}

func NewOrderHandler(orderService ports.OrderService) *OrderHandler {
	return &OrderHandler{srv: orderService}
}

func (oh *OrderHandler) SaveOrder(r *http.Request, w http.ResponseWriter) {
	type parameters struct {
		ID                *uuid.UUID              `json:"id,omitempty"`
		Provider          domain.Providers        `json:"provider"`
		UserID            uuid.UUID               `json:"user_id"`
		PaymentID         *string                 `json:"payment_id,omitempty"`
		SecureToken       *uuid.UUID              `json:"secure_token,omitempty"`
		ExternalReference *string                 `json:"external_reference,omitempty"`
		Currency          domain.Currencies       `json:"currency"`
		SubTotal          float64                 `json:"sub_total"`
		Disscount         *float64                `json:"disscount,omitempty"`
		DisscountType     *domain.DisscountTypes  `json:"disscount_type,omitempty"`
		Total             float64                 `json:"total"`
		Paid              bool                    `json:"paid"`
		PayStatus         domain.PayStatus        `json:"pay_status"`
		PayStatusDetail   *domain.PayStatusDetail `json:"pay_status_detail,omitempty"`
		ExpiresAt         string                  `json:"expires_at"`
	}

	params, err := utils.ParseRequestBody[parameters](r)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
		return
	}

	// Verify HTTP method
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodPut && params.ID == nil {
		http.Error(w, "Method PUT requires the OrderID", http.StatusMethodNotAllowed)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodPost && params.ID != nil {
		http.Error(w, "Method POST doesn't accept a OrderID", http.StatusMethodNotAllowed)
		return
	}

	// Validate required fields
	if params.UserID == uuid.Nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID is required")
		return
	}
	if params.Currency == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "Currency is required")
		return
	}
	if params.SubTotal <= 0 {
		httpdtos.RespondError(w, http.StatusBadRequest, "SubTotal must be greater than zero")
		return
	}
	if params.Disscount != nil && *params.Disscount < 0 {
		httpdtos.RespondError(w, http.StatusBadRequest, "Disscount must be greater than or equal to zero")
		return
	}
	if params.Total <= 0 {
		httpdtos.RespondError(w, http.StatusBadRequest, "Total must be greater than zero")
		return
	}
	if params.PayStatus == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "PayStatus is required")
		return
	}

	// Set id to avoid panic
	var id uuid.UUID
	if params.ID == nil {
		id = uuid.Nil
	} else {
		id = *params.ID
	}

	inputs := ports.SaveOrderInputs{
		ID:                id,
		UserID:            params.UserID,
		PaymentID:         *params.PaymentID,
		ExternalReference: *params.ExternalReference,
		Currency:          params.Currency,
		SubTotal:          params.SubTotal,
		Disscount:         params.Disscount,
		DisscountTypes:    params.DisscountType,
		PayStatus:         params.PayStatus,
		PayStatusDetail:   *params.PayStatusDetail,
	}

	result, err := oh.srv.SaveOrder(r.Context(), inputs)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error saving order: %s", err))
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Order successfully saved", result)
}

func (oh *OrderHandler) GetOrderByID(r *http.Request, w http.ResponseWriter) {
	// Verify HTTP method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract order ID from URL parameters
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "OrderID is required")
		return
	}

	// Parse the order ID to UUID
	parsedOrderId, err := uuid.Parse(orderID)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid OrderID: %s", err))
		return
	}

	// Retrieve the order using the service
	order, err := oh.srv.GetOrderById(r.Context(), parsedOrderId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving order: %s", err))
		return
	}
	if order == nil {
		httpdtos.RespondError(w, http.StatusNotFound, "Order not found")
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Order retrieved successfully", order)
}

func (oh *OrderHandler) GetAllOrders(r *http.Request, w http.ResponseWriter) {
	// Verify HTTP method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the order using the service
	orders, err := oh.srv.ListOrders(r.Context())
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving orders: %s", err))
		return
	}
	if orders == nil {
		httpdtos.RespondError(w, http.StatusNotFound, "Orders not found")
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Orders retrieved successfully", orders)
}
