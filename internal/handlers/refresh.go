package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/mail.v2"

	"github.com/volchok96/auth-medods/internal/database"
	"github.com/volchok96/auth-medods/internal/domain/api/response"
	"github.com/volchok96/auth-medods/internal/domain/ip"
	"github.com/volchok96/auth-medods/internal/domain/jwt"
	"golang.org/x/crypto/bcrypt"
)

func RefreshHandler(db database.DBInterface, ownKey string, tokenTTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp response.RefreshToken
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			log.Error().Err(err).Msg("failed to decode body")
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		log.Info().
			Str("GUID", resp.GUID).
			Str("RefreshToken", resp.RefreshToken).
			Msg("Received refresh token request")

		if len(resp.RefreshToken) == 0 {
			log.Error().Msg("Refresh token is empty")
			http.Error(w, "refresh token is required", http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByGUID(resp.GUID)
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

		// Email notification when IP changes
		if user.IP != clientIP {
			log.Warn().
				Str("User email", user.Email).
				Str("Old IP", user.IP).
				Str("New IP", clientIP).
				Msg("IP address changed")

			emailBody := fmt.Sprintf("Query from a new IP address (%s). Was it you?", clientIP)
			if err := emailWarning(user.Email, emailBody); err != nil {
				log.Error().
					Err(err).
					Msg("failed to send message to email")
			}
		}

		// Decode the token
		log.Info().Str("RefreshToken", resp.RefreshToken).Msg("Decoding refresh token")

		decodedToken, err := base64.StdEncoding.DecodeString(resp.RefreshToken)
		if err != nil {
			log.Error().Err(err).Msg("failed to decode refresh token")
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}

		if len(decodedToken) == 0 {
			log.Error().Msg("decoded token is empty")
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.HashedRefreshToken), decodedToken)
		if err != nil {
			log.Error().Err(err).Msg("invalid refresh token")
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		newAccessToken, newRefreshToken, newRefreshTokenHash, err := jwt.NewTokens(user.UserGUID.String(), ownKey, clientIP, int(tokenTTL.Hours()))
		if err != nil {
			log.Error().Err(err).Msg("failed to generate new tokens")
			http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
			return
		}

		user.IP = clientIP
		user.HashedRefreshToken = newRefreshTokenHash
		if err := db.UpdateUser(user); err != nil {
			log.Error().Err(err).Msg("failed to update user")
			http.Error(w, "failed to update", http.StatusInternalServerError)
			return
		}

		// Convert token to base64
		refreshBase64 := base64.StdEncoding.EncodeToString([]byte(newRefreshToken))
		response := response.UserResponse{
			AccessToken:     newAccessToken,
			GetRefreshToken: refreshBase64,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error().Err(err).Msg("failed to encode response")
		}

		log.Info().
			Str("status", "success").
			Int("code", http.StatusOK).
			Msg("Successfully sent response")
	}
}

func emailWarning(email, bodyString string) error {
	const op = "handlers.refreshHandler.emailWarning"
	m := mail.NewMessage()

	m.SetHeader("From", "your_login@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "WARNING! IP ADDRESS TO ACCESS MEDODS HAS JUST CHANGED")
	m.SetBody("text/plain", bodyString)

	// Use environment variables for login credentials.
	d := mail.NewDialer("smtp.gmail.com", 587, "your_login@gmail.com", "your_password")

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
