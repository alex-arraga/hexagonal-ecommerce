package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce/internal/adapters/api/http/handlers"
	"go-ecommerce/internal/adapters/api/http/routes"
	"go-ecommerce/internal/adapters/security"
	"go-ecommerce/internal/adapters/shared/encoding"
	"go-ecommerce/internal/test_helpers/test_containers"

	"go-ecommerce/internal/adapters/storage/cache/redis"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ! Run docker demon to execute these tests
func Test_UserE2E(t *testing.T) {
	ctx := context.Background()

	// init postgres container
	postgresCont, err := test_containers.NewPostgresContainerDB(t)
	require.NoError(t, err)
	defer postgresCont.Container.Terminate(ctx)

	// init redis container
	redisCont, err := test_containers.NewRedisContainer(t)
	require.NoError(t, err)
	defer redisCont.Container.Terminate(ctx)

	// init transaction
	tx := postgresCont.DB.Begin()
	t.Cleanup(func() { tx.Rollback() })

	hasher := &security.Hasher{}

	// dependency injection
	repo := repository.NewUserRepo(tx)
	srv := services.NewUserService(repo, redisCont.Client, hasher)
	handler := handlers.NewUserHandler(srv)

	r := chi.NewRouter()
	routes.LoadUserRoutes(r, handler)

	server := httptest.NewServer(r)
	defer server.Close()

	// --------------------
	// Step 1 - Create user by POST
	// --------------------
	payload := map[string]string{
		"name":     "John",
		"password": "password",
		"email":    "john@mail.test",
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	serverUrl := fmt.Sprintf("%s/user", server.URL)
	contentType := "application/json"

	// http post
	resp, err := http.Post(serverUrl, contentType, bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]any
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]any) // assert to manage "data" prop in body of response
	userID := data["ID"]

	assert.Equal(t, payload["name"], data["Name"])
	assert.Equal(t, payload["email"], data["Email"])

	// --------------------
	// Step 3 - Validate that caching in redis
	// --------------------
	var user map[string]any

	cacheKey := redis.GenerateCacheKey("user", userID)
	val, err := redisCont.Client.Get(ctx, cacheKey)
	require.NoError(t, err)

	err = encoding.Deserialize(val, &user)
	require.NoError(t, err)

	assert.Contains(t, user["Name"], payload["Name"])
	assert.Contains(t, user["Email"], payload["Email"])
}

/*

// --------------------
TODO Step 2 - Get user by GET
// --------------------
serverUrlWithParam := fmt.Sprintf("%s/user/%s", server.URL, userID)

getResp, err := http.Get(serverUrlWithParam)
require.NoError(t, err)

assert.Equal(t, http.StatusOK, getResp.StatusCode)

// read body
var fetchedUser map[string]any

bodyBytes, err := io.ReadAll(getResp.Body)
require.NoError(t, err)
json.Unmarshal(bodyBytes, &fetchedUser)

assert.Equal(t, userID, fetchedUser["id"])
assert.Equal(t, payload["name"], fetchedUser["Name"])
assert.Equal(t, payload["email"], fetchedUser["Email"])
*/
