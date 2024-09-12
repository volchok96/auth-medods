package pgsql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/volchok96/auth-medods/internal/database/models"
)

type DB struct {
	db *sql.DB
}

func NewDB(connStr string) (*DB, error) {
	const fn = "database.pgsql.NewDB"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) UpdateUser(user *models.User) error {
	query := `
		INSERT INTO users (user_guid, ip, hashed_refresh_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_guid) 
		DO UPDATE SET ip = $2, hashed_refresh_token = $3
	`

	_, err := db.db.Exec(query, user.UserGUID, user.IP, user.HashedRefreshToken)

	return err
}

func (db *DB) GetUserByGUID(guid string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, user_guid, ip, hashed_refresh_token 
		  FROM users 
		  WHERE user_guid = $1`

	err := db.db.QueryRow(query, guid).Scan(&user.ID, &user.UserGUID, &user.IP, &user.HashedRefreshToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}
