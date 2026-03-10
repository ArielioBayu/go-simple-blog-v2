package memberships

import (
	"time"
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RefreshTokenRequest struct {
		Token string `json:"token"`
	}
)

type UserModel struct {
	ID        int64     `column:"id"`
	Email     string    `column:"email"`
	Password  string    `column:"password"`
	CreatedAt time.Time `column:"created_at"`
	UpdatedAt time.Time `column:"updated_at"`
	CreatedBy string    `column:"created_by"`
	UpdatedBy string    `column:"updated_by"`
	Username  string    `column:"username"`
}

type RefreshTokenModel struct {
	ID           int       `json:"id" column:"id"`
	UserId       int       `json:"user_id" column:"user_id"`
	RefreshToken string    `json:"refresh_token" column:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at" column:"expired_at"`
	CreatedAt    time.Time `json:"created_at" column:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" column:"updated_at"`
	CreatedBy    string    `json:"created_by" column:"created_by"`
	UpdatedBy    string    `json:"updated_by" column:"updated_by"`
}
