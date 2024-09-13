package unit_tests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/volchok96/auth-medods/internal/database/models"
	"github.com/volchok96/auth-medods/internal/domain/api/response"
	"github.com/volchok96/auth-medods/internal/handlers"
)

type RMockDB struct {
	mock.Mock
}

func (m *RMockDB) Close() error {
	panic("unimplemented")
}

func (m *RMockDB) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *RMockDB) GetUserByGUID(guid string) (*models.User, error) {
	args := m.Called(guid)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func TestRefreshHandler(t *testing.T) {
	ownKey := "test_key"
	tokenTTL := 30 * time.Minute
	mockDB := new(RMockDB)
	handler := handlers.RefreshHandler(mockDB, ownKey, tokenTTL)

	t.Run("successful token refresh", func(t *testing.T) {
		guid := uuid.New().String()
		refreshToken := uuid.New().String()
		hashedToken, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

		mockDB.On("GetUserByGUID", guid).Return(&models.User{
			UserGUID:           uuid.MustParse(guid),
			HashedRefreshToken: string(hashedToken),
		}, nil)
		mockDB.On("UpdateUser", mock.Anything).Return(nil)

		body, _ := json.Marshal(response.RefreshToken{
			GUID:         guid,
			RefreshToken: base64.StdEncoding.EncodeToString([]byte(refreshToken)),
		})

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody response.UserResponse
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			t.Fatal(err)
		}

		assert.NotEmpty(t, responseBody.AccessToken)
		assert.NotEmpty(t, responseBody.GetRefreshToken)
	})
}
