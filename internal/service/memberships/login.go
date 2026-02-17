package memberships

import (
	"context"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	repo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type MembershipsService interface {
	GetUser(ctx context.Context, request memberships.SignUpRequest) (*memberships.UserModel, error)
	SignUp(ctx context.Context, request memberships.SignUpRequest) error
	SignIn(ctx context.Context, request memberships.SignInRequest) (string, error)
}

type membershipsService struct {
	cfg             *configs.Config
	membershipsRepo repo.MembershipRepository
}

func NewMembershipsService(cfg *configs.Config, membershipsRepo repo.MembershipRepository) MembershipsService {
	return &membershipsService{cfg, membershipsRepo}
}

func (s *membershipsService) GetUser(ctx context.Context, request memberships.SignUpRequest) (*memberships.UserModel, error) {
	data, err := s.membershipsRepo.GetUser(ctx, request.Email, request.Username)
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

	//	Jika data user belum ada maka jalankan decrypt pwd
	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// deklarasi waktu real time untuk created_at
	time := time.Now()
	model := memberships.UserModel{
		Email:     request.Email,
		Username:  request.Username,
		Password:  string(pass), //	masukkan password yg telah dilakukan bcrypt
		CreatedAt: time,
		UpdatedAt: time,
		CreatedBy: request.Email,
		UpdatedBy: request.Email,
	}

	//	last step, create user
	err = s.membershipsRepo.CreateUser(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (s *membershipsService) SignIn(ctx context.Context, request memberships.SignInRequest) (string, error) {
	user, err := s.membershipsRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", constants.ErrDataNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", constants.ErrInvalidPassword
	}

	token, err := jwt.CreateToken(int(user.ID), user.Username, s.cfg.Service.SecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
