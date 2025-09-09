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

	// Validating if email format is valid
	if isValid := utils.IsValidEmail(params.Email); !isValid {
		httpdtos.RespondError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate len of password
	if len(params.Name) < 3 {
		httpdtos.RespondError(w, http.StatusBadRequest, "The name must contain at least 3 characters")
		return
	}

	// Validate len of password
	if len(params.Password) < 6 {
		httpdtos.RespondError(w, http.StatusBadRequest, "The password must contain at least 6 characters")
		return
	}

	u := &domain.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	}

	user, err := h.us.Register(r.Context(), u)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	httpdtos.RespondJSON(w, http.StatusCreated, "User successfully registered", user)
}
