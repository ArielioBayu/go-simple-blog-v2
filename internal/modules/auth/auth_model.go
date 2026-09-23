package auth

import "time"

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

	VerifyOTPRequest struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	ResendOTPRequest struct {
		Email string `json:"email"`
	}
)

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

type UserOTPModel struct {
	ID        int64     `json:"id" column:"id"`
	UserID    int64     `json:"user_id" column:"user_id"`
	OTPCode   string    `json:"otp_code" column:"otp_code"`
	OTPType   string    `json:"otp_type" column:"otp_type"`
	ExpiredAt time.Time `json:"expired_at" column:"expired_at"`
	CreatedAt time.Time `json:"created_at" column:"created_at"`
	UpdatedAt time.Time `json:"updated_at" column:"updated_at"`
}
