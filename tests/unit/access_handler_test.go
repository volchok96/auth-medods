package unit_tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/volchok96/auth-medods/internal/database/models"
	"github.com/volchok96/auth-medods/internal/handlers"
	"golang.org/x/crypto/bcrypt"
)

type AMockDB struct {
	mock.Mock
}

func (m *AMockDB) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func setupTestData(db *sql.DB, guid string) error {
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte("someRefreshToken"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate hashed token: %w", err)
	}

	query := `INSERT INTO users (user_guid, email, hashed_refresh_token, ip) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(query, guid, "testuser@example.com", hashedRefreshToken, "127.0.0.1")
	if err != nil {
		return fmt.Errorf("failed to insert test data: %w", err)
	}

	return nil
}

func teardownTestData(db *sql.DB, guid string) error {
	query := `DELETE FROM users WHERE user_guid = $1`
	_, err := db.Exec(query, guid)
	if err != nil {
		return fmt.Errorf("failed to delete test data: %w", err)
	}

	return nil
}

func (m *AMockDB) GetUserByGUID(guid string) (*models.User, error) {
	args := m.Called(guid)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AMockDB) Close() error {
	return nil
}

func TestAccessHandler(t *testing.T) {
	ownKey := "test_key"
	tokenTTL := 30 * time.Minute

	db, err := sql.Open("postgres", "postgres://postgres:mypass@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	guid := uuid.New().String()

	err = setupTestData(db, guid)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}

	defer func() {
		err = teardownTestData(db, guid)
		if err != nil {
			t.Fatalf("failed to teardown test data: %v", err)
		}
	}()

	mockDB := new(AMockDB)
	handler := handlers.AccessHandler(mockDB, ownKey, tokenTTL)

	t.Run("successful access token generation", func(t *testing.T) {
		mockDB.On("GetUserByGUID", guid).Return(&models.User{
			UserGUID:           uuid.MustParse(guid),
			HashedRefreshToken: "someHashedToken",
			IP:                 "127.0.0.1",
			Email:              "testuser@example.com",
		}, nil)

		mockDB.On("UpdateUser", mock.Anything).Return(nil)

		req := httptest.NewRequest(http.MethodGet, "/access?guid="+guid, nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, responseBody, "access_token")
		assert.Contains(t, responseBody, "refresh_token")
	})
}

