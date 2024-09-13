package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/volchok96/auth-medods/internal/database/models"
	"github.com/volchok96/auth-medods/internal/domain/api/response"
	"github.com/volchok96/auth-medods/internal/handlers"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetUserByGUID(guid string) (*models.User, error) {
	args := m.Called(guid)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDB) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) Close() error {
	return nil
}

func TestAccessHandlerIntegration(t *testing.T) {
	ownKey := "test_key"
	tokenTTL := 30 * time.Minute

	// Создаем моковую базу данных
	mockDB := new(MockDB)

	guid := uuid.New().String()

	t.Logf("Generated GUID: %s", guid)

	mockDB.On("GetUserByGUID", mock.Anything).Return(&models.User{
		UserGUID:           uuid.MustParse(guid),
		HashedRefreshToken: "someHashedToken",
		IP:                 "127.0.0.1",
		Email:              "testuser@example.com",
	}, nil).Once() // Ожидаем один вызов

	mockDB.On("UpdateUser", mock.Anything).Return(nil).Once() // Ожидаем один вызов

	handler := handlers.AccessHandler(mockDB, ownKey, tokenTTL)

	t.Run("integration test for access handler with mocked DB", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/access?guid="+guid, nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody response.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			t.Fatal("Failed to decode the response body:", err)
		}

		assert.NotEmpty(t, responseBody.AccessToken, "AccessToken should not be empty")
		assert.NotEmpty(t, responseBody.GetRefreshToken, "GetRefreshToken should not be empty")
	})
}
