package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/volchok96/auth-medods/internal/database/pgsql"
	"github.com/volchok96/auth-medods/internal/domain/api/response"
	"github.com/volchok96/auth-medods/internal/domain/ip"
	"github.com/volchok96/auth-medods/internal/domain/jwt"
	"golang.org/x/crypto/bcrypt"
)

func RefreshHandler(pgsql *pgsql.DB, ownKey string, tokenTTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp response.RefreshToken
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			log.Error().Err(err).Msg("failed to decode body")
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := pgsql.GetUserByGUID(resp.GUID)
		if err != nil {
			log.Error().Err(err).Msg("user not found or invalid refresh token")
			http.Error(w, "permission denied", http.StatusUnauthorized)
			return
		}

		clientIP := ip.GetIp(r)
		if clientIP == "" {
			log.Error().Msg("failed to get IP")
			http.Error(w, "failed to get IP", http.StatusInternalServerError)
			return
		}

		decodedToken, err := base64.StdEncoding.DecodeString(resp.RefreshToken)
		if err != nil {
			log.Error().Err(err).Msg("failed to decode refresh token")
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HashedRefreshToken), decodedToken); err != nil {
			log.Error().Err(err).Msg("invalid refresh token")
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		newAccessToken, err := jwt.GenerateJWT(user.UserGUID.String(), ownKey, clientIP, int(tokenTTL.Hours()))
		if err != nil {
			log.Error().Err(err).Msg("failed to generate new access token")
			http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
			return
		}

		newRefreshToken, newRefreshTokenHash, err := jwt.GenerateRefreshToken()
		if err != nil {
			log.Error().Err(err).Msg("failed to generate new refresh token")
			http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
			return
		}

		user.IP = clientIP
		user.HashedRefreshToken = newRefreshTokenHash
		if err := pgsql.UpdateUser(user); err != nil {
			log.Error().Err(err).Msg("failed to update user")
			http.Error(w, "failed to update", http.StatusInternalServerError)
			return
		}

		// Token -> base64
		refreshBase64 := base64.StdEncoding.EncodeToString([]byte(newRefreshToken))
		response := response.UserResponse{AccessToken: newAccessToken, RefreshToken: refreshBase64}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

