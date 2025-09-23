package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/adapters/api/http/utils"
	"go-ecommerce/internal/core/ports"
	"net/http"

	"github.com/google/uuid"
)

type CategoryHandler struct {
	srv ports.CategoryService
}

func NewCategoryHandler(srv ports.CategoryService) *CategoryHandler {
	return &CategoryHandler{srv: srv}
}

func (ch *CategoryHandler) SaveCategory(r *http.Request, w http.ResponseWriter) {
	type parameters struct {
		ID   *uuid.UUID `json:"id,omitempty"`
		Name *string    `json:"name"`
	}

	params, err := utils.ParseRequestBody[parameters](r)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
		return
	}

	if *params.Name == "" || params.Name == nil {
		httpdtos.RespondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	category, err := ch.srv.RegisterCategory(r.Context(), *params.Name)
	if err != nil {
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
	}

	// if the user send a update request, respond "ok = 200", else response "created = 201"
	if params.ID != nil {
		httpdtos.RespondJSON(w, http.StatusOK, "Category successfully updated", category)
	} else {
		httpdtos.RespondJSON(w, http.StatusCreated, "Category successfully created", category)
	}
}
