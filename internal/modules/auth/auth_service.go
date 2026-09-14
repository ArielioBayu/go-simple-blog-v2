package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/jwt"
	refToken "github.com/ArielioBayu/go-simple-blog-v2/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx context.Context, request SignUpRequest) error
	SignIn(ctx context.Context, request SignInRequest) (string, string, error)
	GetIdRefreshToken(ctx context.Context, request RefreshTokenRequest) (*RefreshTokenModel, error)
	ValidateRefreshToken(ctx context.Context, userId int, request RefreshTokenRequest) (string, error)
}

type authService struct {
	cfg      *configs.Config
	authRepo AuthRepository
	userRepo user.UserRepository
}

func NewAuthService(cfg *configs.Config, authRepo AuthRepository, userRepo user.UserRepository) AuthService {
	return &authService{
		cfg:      cfg,
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *authService) SignUp(ctx context.Context, request SignUpRequest) error {
	existingUser, err := s.userRepo.GetUser(ctx, request.Email, request.Username)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return constants.ErrUsernameOrEmailAlreadyExists
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	model := user.UserModel{
		Email:     request.Email,
		Username:  request.Username,
		Password:  string(pass),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: request.Email,
		UpdatedBy: request.Email,
	}

	err = s.userRepo.CreateUser(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) SignIn(ctx context.Context, request SignInRequest) (string, string, error) {
	u, err := s.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return "", "", err
	}

	if u == nil {
		return "", "", constants.ErrDataNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(request.Password))
	if err != nil {
		return "", "", constants.ErrInvalidPassword
	}

	token, err := jwt.CreateToken(int(u.ID), u.Username, s.cfg.Service.SecretKey)
	if err != nil {
		return "", "", fmt.Errorf("service signin create token: %w", err)
	}

	now := time.Now()

	if err := s.authRepo.DeleteExpiredRefreshTokens(ctx, int(u.ID), now); err != nil {
		log.Printf("[service signin]: failed to delete expired refresh tokens for user %d: %v", u.ID, err)
	}

	existsRefreshToken, err := s.authRepo.GetRefreshToken(ctx, int(u.ID), now)
	if err != nil {
		return "", "", fmt.Errorf("service signin get refresh token: %w", err)
	}

	if existsRefreshToken != nil {
		return token, existsRefreshToken.RefreshToken, nil
	}

	refreshToken := refToken.GenerateRefreshToken()
	if refreshToken == "" {
		return token, "", errors.New("failed to generate refresh token")
	}

	err = s.authRepo.InsertRefreshToken(ctx, RefreshTokenModel{
		UserId:       int(u.ID),
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    u.Username,
		UpdatedBy:    u.Username,
	})
	if err != nil {
		return "", "", fmt.Errorf("service insert refresh token: %w", err)
	}

	return token, refreshToken, nil
}

func (s *authService) GetIdRefreshToken(ctx context.Context, request RefreshTokenRequest) (*RefreshTokenModel, error) {
	idRefreshToken, err := s.authRepo.GetIdRefreshToken(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("service GetIdRefreshToken: %w", err)
	}

	if idRefreshToken == nil {
		return nil, constants.ErrRefreshTokenNotFound
	}

	return idRefreshToken, nil
}

func (s *authService) ValidateRefreshToken(ctx context.Context, userId int, request RefreshTokenRequest) (string, error) {
	now := time.Now()

	existRefreshToken, err := s.authRepo.GetRefreshToken(ctx, userId, now)
	if err != nil {
		return "", fmt.Errorf("service get refresh token: %w", err)
	}

	if existRefreshToken == nil {
		return "", constants.ErrTokenExpired
	}

	if existRefreshToken.RefreshToken != request.Token {
		return "", constants.ErrInvalidToken
	}

	u, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("service get user by id: %w", err)
	}

	token, err := jwt.CreateToken(int(u.ID), u.Username, s.cfg.Service.SecretKey)
	if err != nil {
		return "", fmt.Errorf("service create token: %w", err)
	}

	return token, nil
}
