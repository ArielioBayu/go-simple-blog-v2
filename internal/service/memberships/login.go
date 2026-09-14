package memberships

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/jwt"
	refToken "github.com/ArielioBayu/go-simple-blog-v2/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

func (s *membershipsService) GetUser(ctx context.Context, request memberships.SignUpRequest) (*memberships.UserModel, error) {
	data, err := s.membershipsRepo.GetUser(ctx, request.Email, request.Username)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *membershipsService) GetUserById(ctx context.Context, id int) (*memberships.UserModel, error) {
	data, err := s.membershipsRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *membershipsService) SignUp(ctx context.Context, request memberships.SignUpRequest) error {
	user, err := s.membershipsRepo.GetUser(ctx, request.Email, request.Username)
	if err != nil {
		return err
	}

	if user != nil {
		return constants.ErrUsernameOrEmailAlreadyExists
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	time := time.Now()
	model := memberships.UserModel{
		Email:     request.Email,
		Username:  request.Username,
		Password:  string(pass),
		CreatedAt: time,
		UpdatedAt: time,
		CreatedBy: request.Email,
		UpdatedBy: request.Email,
	}

	err = s.membershipsRepo.CreateUser(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (s *membershipsService) SignIn(ctx context.Context, request memberships.SignInRequest) (string, string, error) {
	user, err := s.membershipsRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", constants.ErrDataNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", "", constants.ErrInvalidPassword
	}

	token, err := jwt.CreateToken(int(user.ID), user.Username, s.cfg.Service.SecretKey)
	if err != nil {
		return "", "", fmt.Errorf("service signin create token: %w", err)
	}

	now := time.Now()

	if err := s.membershipsRepo.DeleteExpiredRefreshTokens(ctx, int(user.ID), now); err != nil {
		log.Printf("[service signin]: failed to delete expired refresh tokens for user %d: %v", user.ID, err)
	}

	existsRefreshToken, err := s.membershipsRepo.GetRefreshToken(ctx, int(user.ID), now)
	if err != nil {
		return "", "", fmt.Errorf("service signin get refresh token: %w", err)
	}

	if existsRefreshToken != nil {
		return token, existsRefreshToken.RefreshToken, nil
	}
	// Generate refreshToken jika setelah di cek di db belum ada
	refreshToken := refToken.GenerateRefreshToken()
	if refreshToken == "" {
		return token, "", errors.New("failed to generate refresh token")
	}

	//	Insert refreshToken ke db
	err = s.membershipsRepo.InsertRefreshToken(ctx, memberships.RefreshTokenModel{
		UserId:       int(user.ID),
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    user.Username,
		UpdatedBy:    user.Username,
	})
	if err != nil {
		return "", "", fmt.Errorf("service insert refresh token: %w", err)
	}

	return token, refreshToken, nil
}
