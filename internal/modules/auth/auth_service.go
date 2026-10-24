package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/jwt"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/mail"
	refToken "github.com/ArielioBayu/go-simple-blog-v2/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService interface {
	SignUp(ctx context.Context, request SignUpRequest) error
	VerifyOTP(ctx context.Context, request VerifyOTPRequest) error
	ResendOTP(ctx context.Context, request ResendOTPRequest) error
	SignIn(ctx context.Context, request SignInRequest) (string, string, error)
	GetIdRefreshToken(ctx context.Context, request RefreshTokenRequest) (*RefreshTokenModel, error)
	ValidateRefreshToken(ctx context.Context, userId int, request RefreshTokenRequest) (string, error)
	SignOut(ctx context.Context, userId int) error
	GetGoogleAuthURL(state string) string
	HandleGoogleCallback(ctx context.Context, code string) (string, string, error)
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
		if existingUser.IsVerified {
			return constants.ErrUsernameOrEmailAlreadyExists
		}
		otpCode, err := mail.GenerateOTP()
		if err != nil {
			return constants.ErrFailedGenerateOTP
		}
		_ = s.authRepo.DeleteOTPByUserIDAndType(ctx, int(existingUser.ID), constants.OTPTypeEmailVerification)
		now := time.Now()
		err = s.authRepo.InsertOTP(ctx, UserOTPModel{
			UserID:    existingUser.ID,
			OTPCode:   otpCode,
			OTPType:   constants.OTPTypeEmailVerification,
			ExpiredAt: now.Add(constants.OTPExpirationDuration),
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return fmt.Errorf("service signup save otp: %w", err)
		}
		go func() {
			if err := mail.SendOTPEmail(s.cfg.SMTP, existingUser.Email, existingUser.Username, otpCode); err != nil {
				log.Printf("[service signup] warning: failed to send otp email: %v", err)
			}
		}()
		return nil
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	model := user.UserModel{
		Email:      request.Email,
		Username:   request.Username,
		Password:   string(pass),
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
		CreatedBy:  request.Email,
		UpdatedBy:  request.Email,
	}

	err = s.userRepo.CreateUser(ctx, &model)
	if err != nil {
		return err
	}

	otpCode, err := mail.GenerateOTP()
	if err != nil {
		return constants.ErrFailedGenerateOTP
	}

	err = s.authRepo.InsertOTP(ctx, UserOTPModel{
		UserID:    model.ID,
		OTPCode:   otpCode,
		OTPType:   constants.OTPTypeEmailVerification,
		ExpiredAt: now.Add(constants.OTPExpirationDuration),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("service signup save otp: %w", err)
	}

	go func() {
		if err := mail.SendOTPEmail(s.cfg.SMTP, model.Email, model.Username, otpCode); err != nil {
			log.Printf("[service signup] warning: failed to send otp email: %v", err)
		}
	}()

	return nil
}

func (s *authService) VerifyOTP(ctx context.Context, request VerifyOTPRequest) error {
	u, err := s.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if u == nil {
		return constants.ErrUserNotFound
	}

	if u.IsVerified {
		return constants.ErrAccountAlreadyVerified
	}

	validOTP, err := s.authRepo.GetValidOTP(ctx, int(u.ID), request.OTP, constants.OTPTypeEmailVerification, time.Now())
	if err != nil {
		return fmt.Errorf("service verify otp: %w", err)
	}

	if validOTP == nil {
		return constants.ErrInvalidOrExpiredOTP
	}

	err = s.userRepo.UpdateUserVerification(ctx, int(u.ID), true)
	if err != nil {
		return fmt.Errorf("service verify otp update user: %w", err)
	}

	_ = s.authRepo.DeleteOTPByUserIDAndType(ctx, int(u.ID), constants.OTPTypeEmailVerification)

	return nil
}

func (s *authService) ResendOTP(ctx context.Context, request ResendOTPRequest) error {
	u, err := s.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if u == nil {
		return constants.ErrUserNotFound
	}

	if u.IsVerified {
		return constants.ErrAccountAlreadyVerified
	}

	otpCode, err := mail.GenerateOTP()
	if err != nil {
		return constants.ErrFailedGenerateOTP
	}

	_ = s.authRepo.DeleteOTPByUserIDAndType(ctx, int(u.ID), constants.OTPTypeEmailVerification)

	now := time.Now()
	err = s.authRepo.InsertOTP(ctx, UserOTPModel{
		UserID:    u.ID,
		OTPCode:   otpCode,
		OTPType:   constants.OTPTypeEmailVerification,
		ExpiredAt: now.Add(constants.OTPExpirationDuration),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("service resend otp save: %w", err)
	}

	go func() {
		if err := mail.SendOTPEmail(s.cfg.SMTP, u.Email, u.Username, otpCode); err != nil {
			log.Printf("[service resend otp] warning: failed to send otp email: %v", err)
		}
	}()

	return nil
}

func (s *authService) SignIn(ctx context.Context, request SignInRequest) (string, string, error) {
	identifier := strings.TrimSpace(request.Identifier)
	if identifier == "" {
		identifier = strings.TrimSpace(request.Username)
	}
	if identifier == "" {
		identifier = strings.TrimSpace(request.Email)
	}

	if identifier == "" {
		return "", "", constants.ErrDataNotFound
	}

	u, err := s.userRepo.GetUserByEmailOrUsername(ctx, identifier)
	if err != nil {
		return "", "", err
	}

	if u == nil {
		return "", "", constants.ErrDataNotFound
	}

	if u.Password == "" {
		return "", "", constants.ErrInvalidPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(request.Password))
	if err != nil {
		return "", "", constants.ErrInvalidPassword
	}

	if !u.IsVerified {
		return "", "", constants.ErrAccountNotVerified
	}

	return s.generateTokens(ctx, u)
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

func (s *authService) SignOut(ctx context.Context, userId int) error {
	err := s.authRepo.DeleteRefreshTokenByUserId(ctx, userId)
	if err != nil {
		return fmt.Errorf("service SignOut: %w", err)
	}
	return nil
}

func (s *authService) getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.cfg.OAuth.Google.ClientID,
		ClientSecret: s.cfg.OAuth.Google.ClientSecret,
		RedirectURL:  s.cfg.OAuth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (s *authService) GetGoogleAuthURL(state string) string {
	oauthConfig := s.getGoogleOAuthConfig()
	return oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *authService) HandleGoogleCallback(ctx context.Context, code string) (string, string, error) {
	if s.cfg.OAuth.Google.ClientID == "" || s.cfg.OAuth.Google.ClientSecret == "" {
		return "", "", errors.New("google oauth client_id or client_secret is not configured")
	}

	oauthConfig := s.getGoogleOAuthConfig()
	oauthToken, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("failed to exchange code for oauth token: %w", err)
	}

	client := oauthConfig.Client(ctx, oauthToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch user info from google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("google userinfo returned unexpected status: %d", resp.StatusCode)
	}

	var googleUserInfo struct {
		ID            string `json:"id"`
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUserInfo); err != nil {
		return "", "", fmt.Errorf("failed to decode google userinfo: %w", err)
	}

	googleID := googleUserInfo.ID
	if googleID == "" {
		googleID = googleUserInfo.Sub
	}
	if googleID == "" {
		return "", "", errors.New("google account did not provide a valid id")
	}

	userEmail := strings.TrimSpace(strings.ToLower(googleUserInfo.Email))
	if userEmail == "" {
		return "", "", errors.New("google account did not provide a valid email")
	}

	u, err := s.userRepo.GetUserByGoogleID(ctx, googleID)
	if err != nil {
		return "", "", fmt.Errorf("service get user by google id: %w", err)
	}

	if u == nil {
		u, err = s.userRepo.GetUserByEmail(ctx, userEmail)
		if err != nil {
			return "", "", fmt.Errorf("service get user by email: %w", err)
		}

		if u != nil {
			if err := s.userRepo.LinkGoogleAccount(ctx, u.ID, googleID); err != nil {
				return "", "", fmt.Errorf("service link google account: %w", err)
			}
			u.GoogleID = &googleID
			u.IsVerified = true
		} else {
			// 3. New user registration via Google OAuth
			username := s.generateUniqueUsername(ctx, userEmail, googleUserInfo.Name)
			now := time.Now()
			newUser := &user.UserModel{
				Email:      userEmail,
				Username:   username,
				AvatarURL:  googleUserInfo.Picture,
				GoogleID:   &googleID,
				IsVerified: true,
				IsPrivate:  false,
				CreatedAt:  now,
				UpdatedAt:  now,
				CreatedBy:  "google_oauth",
				UpdatedBy:  "google_oauth",
			}

			if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
				return "", "", fmt.Errorf("service create oauth user: %w", err)
			}
			u = newUser
		}
	}

	return s.generateTokens(ctx, u)
}

func (s *authService) generateTokens(ctx context.Context, u *user.UserModel) (string, string, error) {
	token, err := jwt.CreateToken(int(u.ID), u.Username, s.cfg.Service.SecretKey)
	if err != nil {
		return "", "", fmt.Errorf("service create token: %w", err)
	}

	now := time.Now()

	if err := s.authRepo.DeleteExpiredRefreshTokens(ctx, int(u.ID), now); err != nil {
		log.Printf("[service auth]: failed to delete expired refresh tokens for user %d: %v", u.ID, err)
	}

	existsRefreshToken, err := s.authRepo.GetRefreshToken(ctx, int(u.ID), now)
	if err != nil {
		return "", "", fmt.Errorf("service get refresh token: %w", err)
	}

	if existsRefreshToken != nil {
		return token, existsRefreshToken.RefreshToken, nil
	}

	refreshToken := refToken.GenerateRefreshToken()
	if refreshToken == "" {
		return token, "", constants.ErrFailedGenerateRefreshToken
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

func (s *authService) generateUniqueUsername(ctx context.Context, email, name string) string {
	var base string
	if parts := strings.Split(email, "@"); len(parts) > 0 && parts[0] != "" {
		base = parts[0]
	} else if name != "" {
		base = strings.ToLower(strings.ReplaceAll(name, " ", ""))
	} else {
		base = "user"
	}

	var cleaned strings.Builder
	for _, r := range strings.ToLower(base) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			cleaned.WriteRune(r)
		}
	}
	candidate := cleaned.String()
	if candidate == "" {
		candidate = "user"
	}

	for i := 0; i < 5; i++ {
		existing, err := s.userRepo.GetUser(ctx, "", candidate)
		if err == nil && existing == nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s%d", cleaned.String(), rand.Intn(9000)+1000)
	}

	return fmt.Sprintf("%s%d", cleaned.String(), time.Now().Unix()%10000)
}
