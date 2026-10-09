package apperror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
	"github.com/go-playground/validator/v10"
)

// AppError is the standard custom error type for application-wide errors.
// It decouples internal error details from sanitized public responses.
type AppError struct {
	StatusCode int               `json:"status"`
	Message    string            `json:"message"`
	Code       string            `json:"code,omitempty"`
	Errors     map[string]string `json:"errors,omitempty"`
	Err        error             `json:"-"` // Root technical error (for logging only)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Factory Constructors
func New(statusCode int, code string, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func NewBadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, err)
}

func NewValidation(message string, errs map[string]string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    message,
		Errors:     errs,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message, err)
}

func NewForbidden(message string, err error) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message, err)
}

func NewNotFound(message string, err error) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message, err)
}

func NewConflict(message string, err error) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, err)
}

func NewTooManyRequests(message string, err error) *AppError {
	return New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message, err)
}

func NewInternal(message string, err error) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, err)
}

// FromError converts any standard error into an AppError with appropriate status code and sanitized message.
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	// 1. Check if it's already an *AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// 2. Handle Gin / Go Playground Validator errors
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		errMap := make(map[string]string)
		for _, fieldErr := range valErrs {
			field := fieldErr.Field()
			tag := fieldErr.Tag()
			param := fieldErr.Param()

			switch tag {
			case "required":
				errMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errMap[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				errMap[field] = fmt.Sprintf("%s must be at least %s characters", field, param)
			case "max":
				errMap[field] = fmt.Sprintf("%s must not exceed %s characters", field, param)
			default:
				errMap[field] = fmt.Sprintf("%s failed validation on '%s'", field, tag)
			}
		}
		return NewValidation("Validation failed", errMap)
	}

	// 3. Match against known application sentinel errors
	switch {
	// Not Found (404)
	case errors.Is(err, constants.ErrDataNotFound):
		return NewNotFound("Data Not Found", err)
	case errors.Is(err, constants.ErrUserNotFound):
		return NewNotFound("user not found", err)
	case errors.Is(err, constants.ErrPostNotFound):
		return NewNotFound("Post Not Found", err)
	case errors.Is(err, constants.ErrCommentNotFound):
		return NewNotFound("comment not found", err)
	case errors.Is(err, constants.ErrReplyNotFound):
		return NewNotFound("reply not found", err)
	case errors.Is(err, constants.ErrRefreshTokenNotFound):
		return NewNotFound("refresh token not found", err)
	case errors.Is(err, constants.ErrFollowNotFound):
		return NewNotFound("follow relationship not found", err)
	case errors.Is(err, constants.ErrFollowRequestNotFound):
		return NewNotFound("follow request not found", err)

	// Unauthorized (401)
	case errors.Is(err, constants.ErrUnauthorized):
		return NewUnauthorized("Unauthorized", err)
	case errors.Is(err, constants.ErrMissingToken):
		return NewUnauthorized("Missing Token", err)
	case errors.Is(err, constants.ErrInvalidToken):
		return NewUnauthorized("invalid token", err)
	case errors.Is(err, constants.ErrTokenExpired):
		return NewUnauthorized("refresh token has expired", err)

	// Forbidden (403)
	case errors.Is(err, constants.ErrForbidden):
		return NewForbidden("forbidden", err)
	case errors.Is(err, constants.ErrAccountNotVerified):
		return NewForbidden("account is not verified, please verify your email first", err)
	case errors.Is(err, constants.ErrPrivateAccount):
		return NewForbidden("this account is private. follow this account to see their posts", err)

	// Conflict (409)
	case errors.Is(err, constants.ErrUsernameOrEmailAlreadyExists):
		return NewConflict("Username or Email Already Exists!", err)
	case errors.Is(err, constants.ErrAlreadyFollowing):
		return NewConflict("already following this user", err)
	case errors.Is(err, constants.ErrFollowRequestPending):
		return NewConflict("follow request already sent", err)

	// Bad Request (400)
	case errors.Is(err, constants.ErrInvalidPassword):
		return NewBadRequest("invalid password", err)
	case errors.Is(err, constants.ErrInvalidOrExpiredOTP):
		return NewBadRequest("invalid or expired OTP code", err)
	case errors.Is(err, constants.ErrAccountAlreadyVerified):
		return NewBadRequest("account is already verified", err)
	case errors.Is(err, constants.ErrCannotFollowSelf):
		return NewBadRequest("cannot follow yourself", err)
	case errors.Is(err, utils.ErrFileRequired), errors.Is(err, utils.ErrFileTooLarge), errors.Is(err, utils.ErrInvalidFormat):
		return NewBadRequest(err.Error(), err)

	// Rate Limited (429)
	case errors.Is(err, constants.ErrTooManyRequests):
		return NewTooManyRequests("too many requests, please try again later", err)

	// Default fallback: 500 Internal Server Error (Never expose raw internal error to public client)
	default:
		return NewInternal("internal server error", err)
	}
}
