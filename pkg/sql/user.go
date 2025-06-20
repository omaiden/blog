package sql

import (
	"context"
	"github.com/acoshift/pgsql/pgctx"
	"time"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := pgctx.QueryRow(ctx, `
	SELECT *
		FROM users
		WHERE username = $1
		`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	return &user, err
}

func CreateUser(ctx context.Context, user *User) error {
	return pgctx.QueryRow(ctx, `
        INSERT INTO users (id, username, email, password_hash) 
        	VALUES ($1, $2, $3, $4)
    `, user.ID, user.Username, user.Email, user.PasswordHash).Err()
}
