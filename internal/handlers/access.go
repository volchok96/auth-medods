package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/volchok96/auth-medods/internal/database"
	"github.com/volchok96/auth-medods/internal/database/models"
	"github.com/volchok96/auth-medods/internal/domain/api/response"
	"github.com/volchok96/auth-medods/internal/domain/ip"
	"github.com/volchok96/auth-medods/internal/domain/jwt"
)

// AccessHandler - хендлер для получения токенов доступа
func AccessHandler(db database.DBInterface, ownKey string, tokenTTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := r.URL.Query().Get("guid")
		if guid == "" {
			log.Error().Msg("no guid provided")
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		if _, err := uuid.Parse(guid); err != nil {
			log.Error().Err(err).Msg("invalid guid")
			http.Error(w, "invalid guid", http.StatusBadRequest)
			return
		}

		clientIP := ip.GetIp(r)
		if clientIP == "" {
			log.Error().Msg("failed to get IP")
			http.Error(w, "failed to get IP", http.StatusInternalServerError)
			return
		}

		token, refreshToken, hashedRefreshToken, err := jwt.NewTokens(guid, ownKey, clientIP, int(tokenTTL.Hours()))
		if err != nil {
			log.Error().Err(err).Msg("failed to generate access token")
			http.Error(w, "failed to generate tokens", http.StatusBadRequest)
			return
		}

		GUID, err := uuid.Parse(guid)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse guid")
			http.Error(w, "failed to parse guid", http.StatusInternalServerError)
			return
		}

		user := &models.User{
			UserGUID:           GUID,
			IP:                 clientIP,
			HashedRefreshToken: hashedRefreshToken,
		}

		if err := db.UpdateUser(user); err != nil {
			log.Error().Err(err).Msg("failed to save hash")
			http.Error(w, "failed to save data", http.StatusInternalServerError)
			return
		}

		refreshBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))

		response := response.UserResponse{
			AccessToken:     token,
			GetRefreshToken: refreshBase64,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error().Err(err).Msg("failed to encode response")
		}

		log.Info().
			Str("hashed_refresh_token", string(user.HashedRefreshToken)).
			Str("new_refresh_token", refreshBase64).
			Str("status", "success").
			Int("code", http.StatusOK).
			Msg("Successfully sent response")
	}
}
