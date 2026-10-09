package constants

import "errors"

var (

	//	Error
	ErrUsernameOrEmailAlreadyExists = errors.New("Username or Email Already Exists!")
	ErrMissingToken                 = errors.New("Missing Token")
	ErrUnauthorized                 = errors.New("Unauthorized")
	ErrForbidden                    = errors.New("forbidden")

	// Not Found
	ErrDataNotFound         = errors.New("Data Not Found")
	ErrPostNotFound         = errors.New("Error Post Not Found")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrReplyNotFound        = errors.New("reply not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	//	Invalid Token & Verification
	ErrInvalidToken               = errors.New("invalid token")
	ErrInvalidPassword            = errors.New("invalid password")
	ErrTokenExpired               = errors.New("refresh token has expired")
	ErrFailedGenerateRefreshToken = errors.New("failed to generate refresh token")
	ErrFailedGenerateOTP          = errors.New("failed to generate OTP code")
	ErrAccountNotVerified         = errors.New("account is not verified, please verify your email first")
	ErrInvalidOrExpiredOTP        = errors.New("invalid or expired OTP code")
	ErrAccountAlreadyVerified     = errors.New("account is already verified")
	ErrUserNotFound               = errors.New("user not found")
	ErrTooManyRequests            = errors.New("too many requests, please try again later")

	// Follow & Privacy Errors
	ErrCannotFollowSelf       = errors.New("cannot follow yourself")
	ErrAlreadyFollowing       = errors.New("already following this user")
	ErrFollowRequestPending   = errors.New("follow request already sent")
	ErrFollowNotFound         = errors.New("follow relationship not found")
	ErrFollowRequestNotFound  = errors.New("follow request not found")
	ErrPrivateAccount         = errors.New("this account is private. follow this account to see their posts")
)

