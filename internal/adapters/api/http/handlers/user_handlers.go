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

type UserHandler struct {
	us ports.UserService
}

func NewUserHandler(us ports.UserService) *UserHandler {
	return &UserHandler{us: us}
}

func (h *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role,omitempty"`
	}

	var id uuid.UUID
	userId := chi.URLParam(r, "user_id")

	if userId != "" {
		parsed, err := uuid.Parse(userId)
		if err != nil {
			httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		id = parsed
	} else {
		id = uuid.Nil
	}

	// parse params of body
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

	inputs := domain.SaveUserInputs{
		ID:       id,
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
		Role:     role,
	}

	user, err := h.us.SaveUser(r.Context(), inputs)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	httpdtos.RespondJSON(w, http.StatusCreated, "User successfully registered", user)
}
