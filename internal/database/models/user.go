package models

import "github.com/google/uuid"

type User struct {
	ID               int
	UserGUID         uuid.UUID
	IP           string
	HashedRefreshToken string
	Email string
}