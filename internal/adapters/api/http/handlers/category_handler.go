package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/adapters/api/http/utils"
	"go-ecommerce/internal/core/domain"
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
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
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
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid input: %s", err))
		return
	}

	if params.Name == nil || *params.Name == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "name is required")
		return
	}

	category, err := ch.srv.SaveCategory(r.Context(), uintId, *params.Name)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// if the user send a update request, respond "ok = 200", else response "created = 201"
	if categoryId != "" {
		httpdtos.RespondJSON(w, http.StatusOK, "category successfully updated", category)
		return
	} else {
		httpdtos.RespondJSON(w, http.StatusCreated, "category successfully created", category)
		return
	}
}

func (ch *CategoryHandler) FindCategoryById(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
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
		if err == domain.ErrCategoryNotFound || categ == nil {
			httpdtos.RespondJSON(w, http.StatusNotFound, err.Error(), nil)
			return
		}

		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "category successfully retrieved", categ)
}

func (ch *CategoryHandler) ListCategories(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find all category calling service
	categs, err := ch.srv.ListCategories(r.Context())
	if err != nil {
		if err == domain.ErrCategoryNotFound || len(categs) <= 0 {
			httpdtos.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "categories successfully retrieved", categs)
}

func (ch *CategoryHandler) DeleteCategory(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodDelete {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find id from url param
	categoryId := chi.URLParam(r, "category_id")
	if categoryId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "categoryID is required to delete a category")
		return
	}

	// parse string to uint
	uintId, err := strconv.ParseUint(categoryId, 10, 64)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ch.srv.DeleteCategory(r.Context(), uintId)
	if err != nil {
		if err == domain.ErrCategoryNotFound {
			httpdtos.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "category successfully deleted", nil)
}
