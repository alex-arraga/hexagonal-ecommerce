package handlers

import (
	"fmt"
	httpdtos "go-ecommerce/internal/adapters/api/http/http_dtos"
	"go-ecommerce/internal/adapters/api/http/utils"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"net/http"
)

type UserHandler struct {
	us ports.UserService
}

func NewUserHandler(us ports.UserService) *UserHandler {
	return &UserHandler{us: us}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role,omitempty"`
	}

	params, err := utils.ParseRequestBody[parameters](r)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
		return
	}

	if params.Name == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if params.Email == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if params.Password == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "Password is required")
		return
	}

	role := domain.Client
	if params.Role != "" {
		role = domain.UserRole(params.Role)
	}

	// Validating if email format is valid
	if isValid := utils.IsValidEmail(params.Email); !isValid {
		httpdtos.RespondError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	user, err := h.us.Register(r.Context(), params.Name, params.Email, params.Password, role)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	httpdtos.RespondJSON(w, http.StatusCreated, "User successfully registered", user)
}
