package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/adapters/api/http/utils"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CartHandler struct {
	srv ports.CartService
}

// TODO -> Probar todas las request del carrito, el resto funciona correctamente
func NewCartHandler(srv ports.CartService) *CartHandler {
	return &CartHandler{srv: srv}
}

func (ch *CartHandler) AddProductToCart(r *http.Request, w http.ResponseWriter) {
	// Validate http methods
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	type parameters struct {
		Quantity *int16 `json:"quantity"`
	}

	params, err := utils.ParseRequestBody[parameters](r)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
		return
	}

	// Validate params
	if params.Quantity == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "Quantity is required")
		return
	}

	// Retrieve and validate URL params
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID is required")
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID must be a valid UUID")
		return
	}

	productId := chi.URLParam(r, "product_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID is required")
		return
	}

	parsedProductId, err := uuid.Parse(productId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID must be a valid UUID")
		return
	}

	// Call service to add a product
	err = ch.srv.AddItemToCart(r.Context(), parsedUserId, parsedProductId, *params.Quantity)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding item to cart: %s", err.Error()))
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Product added to cart", nil)
}

func (ch *CartHandler) GetProductsFromCart(r *http.Request, w http.ResponseWriter) {
	// Validate http methods
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Retrieve and validate URL params
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID is required")
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID must be a valid UUID")
		return
	}

	// Call service to add a product
	cart, err := ch.srv.GetCart(r.Context(), parsedUserId)
	if err != nil {
		slog.Error("Error retrieving cart", "user_id", userId, "error", err)
		httpdtos.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving cart: %s", err.Error()))
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Cart successfully retrieved", cart)
}

func (ch *CartHandler) ClearCart(r *http.Request, w http.ResponseWriter) {
	// Validate http methods
	if r.Method != http.MethodDelete {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Retrieve and validate URL params
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID is required")
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID must be a valid UUID")
		return
	}

	// Call service to add a product
	err = ch.srv.Clear(r.Context(), parsedUserId)
	if err != nil {
		if err == domain.ErrAlreadyEmptyCart {
			httpdtos.RespondError(w, http.StatusNoContent, fmt.Sprintf("Error deleting product in cart: %s", err.Error()))
			return
		}

		slog.Error("Error clearing cart", "user_id", userId, "error", err)
		httpdtos.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error clearing cart: %s", err.Error()))
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Cart successfully cleared", nil)
}

func (ch *CartHandler) RemoveItemFromCart(r *http.Request, w http.ResponseWriter) {
	// Validate http methods
	if r.Method != http.MethodDelete {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Retrieve and validate URL params
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID is required")
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "UserID must be a valid UUID")
		return
	}

	productId := chi.URLParam(r, "product_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "ProductID is required")
		return
	}

	parsedProductId, err := uuid.Parse(productId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "ProductID must be a valid UUID")
		return
	}

	// Call service to add a product
	err = ch.srv.RemoveItem(r.Context(), parsedUserId, parsedProductId)
	if err != nil {
		if err == domain.ErrAlreadyEmptyCart {
			httpdtos.RespondError(w, http.StatusNoContent, fmt.Sprintf("Error deleting product in cart: %s", err.Error()))
			return
		}

		if err == domain.ErrProductNotFoundCart {
			httpdtos.RespondError(w, http.StatusNotFound, fmt.Sprintf("Error deleting product in cart: %s", err.Error()))
			return
		}

		slog.Error("Error remove item from cart", "user_id", userId, "product_id", productId, "error", err)
		httpdtos.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting product in cart: %s", err.Error()))
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Cart item successfully removed", nil)
}
