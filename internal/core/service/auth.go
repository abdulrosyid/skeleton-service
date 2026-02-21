package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	repoReg "skeleton-service/internal/adapter/outbound/repository/mysqldb"
	"skeleton-service/internal/core/domain"
	iservice "skeleton-service/internal/port/inbound/service"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrEmailExists   = errors.New("email already registered")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrUserNotFound  = errors.New("user not found")
	ErrInternalError = errors.New("internal error")
)

type AuthService struct {
	repoRegistry     repoReg.RepoSQL
	jwtSecret        string
	jwtExpireMinutes int
}

func NewAuthService(r repoReg.RepoSQL, jwtSecret string, jwtExpireMinutes int) iservice.AuthService {
	return &AuthService{
		repoRegistry:     r,
		jwtSecret:        jwtSecret,
		jwtExpireMinutes: jwtExpireMinutes,
	}
}

func (s *AuthService) Register(ctx context.Context, req iservice.RegisterRequest) (*iservice.UserResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)

	if req.FullName == "" || req.Email == "" || len(req.Password) < 6 {
		return nil, ErrInvalidInput
	}

	repo := s.repoRegistry.GetUserMysqlDbRepository()
	existing, err := repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInternalError
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternalError
	}

	user := &domain.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hashed),
	}

	if err := repo.Create(ctx, user); err != nil {
		return nil, ErrInternalError
	}

	return &iservice.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req iservice.LoginRequest) (*iservice.LoginResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	repo := s.repoRegistry.GetUserMysqlDbRepository()
	user, err := repo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	expiresAt := time.Now().Add(time.Duration(s.jwtExpireMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   expiresAt.Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, ErrInternalError
	}

	return &iservice.LoginResponse{
		Token: token,
		User: iservice.UserResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) Profile(ctx context.Context, userID int64) (*iservice.UserResponse, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}

	repo := s.repoRegistry.GetUserMysqlDbRepository()
	user, err := repo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrInternalError
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &iservice.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}
