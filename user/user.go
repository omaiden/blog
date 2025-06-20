package user

import (
	"blog/internal/auth"
	"blog/pkg/api"
	"context"
	"errors"
	"github.com/moonrhythm/randid"
	"strings"
	"time"

	"blog/pkg/sql"

	dbsql "database/sql"
	"github.com/moonrhythm/validator"
)

type RegisterRequest struct {
	Username     string
	Email        string
	PasswordHash string `json:"password_hash"`
}

func (p *RegisterRequest) Valid() error {
	p.Username = strings.TrimSpace(p.Username)
	p.Email = strings.TrimSpace(p.Email)

	v := validator.New()
	v.Must(p.Username != "", "username required")
	v.Must(p.Email != "", "email required")
	v.Must(p.PasswordHash != "", "password required")
	return v.Error()
}

type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt string
}

func Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	if err := req.Valid(); err != nil {
		return nil, err
	}

	// Check if username exists
	_, err := sql.GetUserByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, dbsql.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		return nil, api.ErrUsernameExist
	}

	// Save user
	sqlUser := sql.User{
		ID:           randid.MustGenerate().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		CreatedAt:    time.Now(),
	}
	err = sql.CreateUser(ctx, &sqlUser)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        sqlUser.ID,
		Username:  sqlUser.Username,
		Email:     sqlUser.Email,
		CreatedAt: api.ConvertTimeToStr(sqlUser.CreatedAt),
	}, nil
}

type LoginRequest struct {
	Username     string
	PasswordHash string `json:"password_hash"`
}

type LoginResponse struct {
	Token string
}

func (p *LoginRequest) Valid() error {
	v := validator.New()
	v.Must(p.Username != "", "username required")
	v.Must(p.PasswordHash != "", "password required")
	return v.Error()
}

func Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if err := req.Valid(); err != nil {
		return nil, err
	}

	user, err := sql.GetUserByUsername(ctx, req.Username)
	if err != nil && !errors.Is(dbsql.ErrNoRows, err) {
		return nil, err
	}

	if errors.Is(dbsql.ErrNoRows, err) || user.PasswordHash != req.PasswordHash {
		return nil, api.ErrLoginFailed
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token}, nil
}
