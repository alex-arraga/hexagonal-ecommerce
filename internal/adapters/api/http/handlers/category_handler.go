package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/adapters/api/http/utils"
	"go-ecommerce/internal/core/ports"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	srv ports.CategoryService
}

func NewCategoryHandler(srv ports.CategoryService) *CategoryHandler {
	return &CategoryHandler{srv: srv}
}

func (ch *CategoryHandler) SaveCategory(r *http.Request, w http.ResponseWriter) {
	type parameters struct {
		Name *string `json:"name"`
	}

	// verify http methods
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// find id from url param
	categoryId := chi.URLParam(r, "category_id")
	var uintId uint64

	if categoryId != "" {
		parsed, err := strconv.ParseUint(categoryId, 10, 64)
		if err != nil {
			httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		uintId = parsed
	}

	// request validations
	params, err := utils.ParseRequestBody[parameters](r)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
		return
	}

	if params.Name == nil || *params.Name == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	category, err := ch.srv.SaveCategory(r.Context(), uintId, *params.Name)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// if the user send a update request, respond "ok = 200", else response "created = 201"
	if categoryId != "" {
		httpdtos.RespondJSON(w, http.StatusOK, "Category successfully updated", category)
	} else {
		httpdtos.RespondJSON(w, http.StatusCreated, "Category successfully created", category)
	}
}

func (ch *CategoryHandler) FindCategoryById(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// find id from url param
	categoryId := chi.URLParam(r, "category_id")
	var uintId uint64

	if categoryId != "" {
		parsed, err := strconv.ParseUint(categoryId, 10, 64)
		if err != nil {
			httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		uintId = parsed
	}

	// find category calling service
	categ, err := ch.srv.GetCategoryByID(r.Context(), uintId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Category successfully retrieved", categ)
}

func (ch *CategoryHandler) ListCategories(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// find all category calling service
	categs, err := ch.srv.ListCategories(r.Context())
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
	}

	httpdtos.RespondJSON(w, http.StatusOK, "Categories successfully retrieved", categs)
}
