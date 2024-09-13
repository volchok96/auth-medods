package database

import "github.com/volchok96/auth-medods/internal/database/models"

type DBInterface interface {
	UpdateUser(user *models.User) error
	GetUserByGUID(guid string) (*models.User, error)
	Close() error
}
