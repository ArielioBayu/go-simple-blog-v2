package memberships

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/jwt"
)

func (s *membershipsService) GetIdRefreshToken(ctx context.Context, request memberships.RefreshTokenRequest) (*memberships.RefreshTokenModel, error) {
	idRefreshToken, err := s.membershipsRepo.GetIdRefreshToken(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("service GetIdRefreshToken: %w", err)
	}

	return idRefreshToken, nil
}

func (s *membershipsService) ValidateRefreshToken(ctx context.Context, userId int, request memberships.RefreshTokenRequest) (string, error) {
	now := time.Now()

	existRefreshToken, err := s.membershipsRepo.GetRefreshToken(ctx, userId, now)
	if err != nil {
		return "", fmt.Errorf("service get refresh token: %w", err)
	}

	if existRefreshToken == nil {
		return "", errors.New("refresh token not found")
	}

	//	cek kalo refresh token di db tidak sesuai sama token yang dimasukkan di client
	if existRefreshToken.RefreshToken != request.Token {
		return "", errors.New("invalid refresh token")
	}

	user, err := s.membershipsRepo.GetUserById(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("service get user by id: %w", err)
	}

	// jika ternyata refresh token sama maka generate jwt untuk mendapatkan access token baru
	token, err := jwt.CreateToken(int(user.ID), user.Username, s.cfg.Service.SecretKey)
	if err != nil {
		return "", fmt.Errorf("service create token: %w", err)
	}

	return token, nil
}
