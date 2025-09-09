package handlers_test

import (
	"context"
	"go-ecommerce/internal/adapters/api/http/handlers"
	"go-ecommerce/internal/adapters/api/http/routes"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/test_helpers/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UserHandler_Register(t *testing.T) {
	mockUserService := &mocks.MockUserService{
		RegisterFunc: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			return &domain.User{
				ID:        uuid.New(),
				Name:      "John",
				Email:     "john@mail.test",
				Password:  "password",
				Role:      domain.Client,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	r := chi.NewRouter()
	handler := handlers.NewUserHandler(mockUserService)
	routes.LoadUserRoutes(r, handler)

	reqBody := `{"name":"John", "email":"john@mail.test", "password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.Register(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
