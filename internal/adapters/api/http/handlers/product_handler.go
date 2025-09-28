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

type ProductHandler struct {
	srv ports.ProductService
}

func NewProductHandler(srv ports.ProductService) *ProductHandler {
	return &ProductHandler{srv: srv}
}

func (ph *ProductHandler) SaveProduct(r *http.Request, w http.ResponseWriter) {
	type parameters struct {
		Name       *string  `json:"name"`
		Image      *string  `json:"image"`
		SKU        *string  `json:"sku"`
		Price      *float64 `json:"price"`
		Stock      *int64   `json:"stock"`
		CategoryID *uint64  `json:"productId"`
	}

	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// extract id from url param
	productId := chi.URLParam(r, "product_id")

	// set id to avoid panic
	var id uuid.UUID
	if productId != "" {
		parsed, err := uuid.Parse(productId)
		if err != nil {
			httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("product id must be a valid uuid: %s", err))
			return
		}
		id = parsed
	} else {
		id = uuid.Nil
	}

	// parse and validates fields of request body
	params, err := utils.ParseRequestBody[parameters](r)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
		return
	}

	if *params.Name == "" || params.Name == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if *params.Image == "" || params.Image == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "image is required")
		return
	}
	if *params.SKU == "" || params.SKU == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "sku is required")
		return
	}
	if params.Price == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "price is required")
		return
	}
	if params.Stock == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "stock is required")
		return
	}
	if params.CategoryID == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "product_id is required")
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
	if id != uuid.Nil {
		httpdtos.RespondJSON(w, http.StatusOK, "product successfully updated", product)
		return
	} else {
		httpdtos.RespondJSON(w, http.StatusCreated, "product successfully created", product)
		return
	}
}

func (ph *ProductHandler) FindProductById(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find id from url param
	productId := chi.URLParam(r, "product_id")
	if productId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "productID is required")
		return
	}

	// convert string id to uuid
	parsedId, err := uuid.Parse(productId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// find category calling service
	prod, err := ph.srv.GetProductById(r.Context(), parsedId)
	if err != nil {
		if err == domain.ErrProductNotFound || prod == nil {
			httpdtos.RespondJSON(w, http.StatusNotFound, err.Error(), nil)
			return
		}

		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "product successfully retrieved", prod)
}

func (ph *ProductHandler) ListProducts(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find all category calling service
	prods, err := ph.srv.ListProducts(r.Context())
	if err != nil {
		if err == domain.ErrProductNotFound || len(prods) <= 0 {
			httpdtos.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "products successfully retrieved", prods)
}

func (ph *ProductHandler) DeleteProduct(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodDelete {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find id from url param
	productId := chi.URLParam(r, "product_id")
	if productId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "productID is required to delete a product")
		return
	}

	// parse string to uint
	parsedId, err := uuid.Parse(productId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ph.srv.DeleteProduct(r.Context(), parsedId)
	if err != nil {
		if err == domain.ErrProductNotFound {
			httpdtos.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "product successfully deleted", nil)
}
