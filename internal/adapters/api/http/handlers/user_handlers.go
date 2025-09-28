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
	"github.com/google/uuid"
)

type UserHandler struct {
	srv ports.UserService
}

func NewUserHandler(us ports.UserService) *UserHandler {
	return &UserHandler{srv: us}
}

func (uh *UserHandler) SaveUser(r *http.Request, w http.ResponseWriter) {
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

	user, err := uh.srv.SaveUser(r.Context(), inputs)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	httpdtos.RespondJSON(w, http.StatusCreated, "User successfully registered", user)
}

func (uh *UserHandler) FindUserById(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find id from url param
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "userID is required")
		return
	}

	// convert string id to uuid
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// find category calling service
	user, err := uh.srv.GetUserByID(r.Context(), parsedId)
	if err != nil {
		if err == domain.ErrUserNotFound || user == nil {
			httpdtos.RespondJSON(w, http.StatusNotFound, err.Error(), nil)
			return
		}

		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "user successfully retrieved", user)
}

func (uh *UserHandler) FindUserByEmail(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find id from url param
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "user email is required")
		return
	}

	// find category calling service
	user, err := uh.srv.GetUserByEmail(r.Context(), userEmail)
	if err != nil {
		if err == domain.ErrUserNotFound || user == nil {
			httpdtos.RespondJSON(w, http.StatusNotFound, err.Error(), nil)
			return
		}

		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "user successfully retrieved", user)
}

func (uh *UserHandler) ListUsers(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodGet {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	skipStr := r.URL.Query().Get("skip")
	limitStr := r.URL.Query().Get("limit")

	skip := 0
	limit := 20

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	if skipStr != "" {
		s, err := strconv.Atoi(skipStr)
		if err == nil && s >= 0 {
			skip = s
		}
	}

	// find all category calling service
	users, err := uh.srv.ListUsers(r.Context(), uint64(skip), uint64(limit))
	if err != nil {
		if err == domain.ErrUserNotFound || len(users) <= 0 {
			httpdtos.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "users successfully retrieved", users)
}

func (uh *UserHandler) DeleteUser(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodDelete {
		httpdtos.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// find id from url param
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		httpdtos.RespondError(w, http.StatusBadRequest, "userID is required to delete a product")
		return
	}

	// parse string to uint
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		httpdtos.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = uh.srv.DeleteUser(r.Context(), parsedId)
	if err != nil {
		if err == domain.ErrUserNotFound {
			httpdtos.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		httpdtos.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpdtos.RespondJSON(w, http.StatusOK, "user successfully deleted", nil)
}
