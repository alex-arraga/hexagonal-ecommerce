package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/adapters/api/http/utils"
	"go-ecommerce/internal/core/ports"
	"net/http"

	"github.com/google/uuid"
)

type ProductHandler struct {
	srv ports.ProductService
}

func NewProductHandler(srv ports.ProductService) *ProductHandler {
	return &ProductHandler{srv: srv}
}

func (ph *ProductHandler) SaveProduct(r *http.Request, w http.ResponseWriter) {
	type parameters struct {
		ID         *uuid.UUID `json:"id,omitempty"`
		Name       *string    `json:"name"`
		Image      *string    `json:"image"`
		SKU        *string    `json:"sku"`
		Price      *float64   `json:"price"`
		Stock      *int64     `json:"stock"`
		CategoryID *uint64    `json:"categoryId"`
	}

	params, err := utils.ParseRequestBody[parameters](r)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
		return
	}

	// set id to avoid panic
	var id uuid.UUID
	if params.ID == nil {
		id = uuid.Nil
	} else {
		id = *params.ID
	}

	if *params.Name == "" || params.Name == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if *params.Image == "" || params.Image == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "Image is required")
		return
	}
	if *params.SKU == "" || params.SKU == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "SKU is required")
		return
	}
	if params.Price == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "Price is required")
		return
	}
	if params.Stock == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "Stock is required")
		return
	}
	if params.CategoryID == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "CategoryID is required")
		return
	}

	// mapping the parameters with inputs required for save a product
	inputs := ports.SaveProductInputs{
		ID:         id,
		Name:       *params.Name,
		Image:      *params.Image,
		SKU:        *params.SKU,
		Price:      *params.Price,
		Stock:      *params.Stock,
		CategoryID: *params.CategoryID,
	}

	product, err := ph.srv.SaveProduct(r.Context(), inputs)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// if the user send a update request, respond "ok = 200", else response "created = 201"
	if params.ID != nil {
		httpdtos.RespondJSON(w, http.StatusOK, "Product successfully updated", product)
		return
	} else {
		httpdtos.RespondJSON(w, http.StatusCreated, "Product successfully created", product)
		return
	}
}
