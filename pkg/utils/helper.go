package utils

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JsonTime time.Time

const (
	layout                   = "2006-01-02 15:04:05"
	OTPTypeEmailVerification = "EMAIL_VERIFICATION"
	OTPExpirationDuration    = 5 * time.Minute
	DefaultUploadDir         = "./uploads"
	UploadDir                = "./uploads"
	MaxUploadSize            = 5 * 1024 * 1024 // 5MB
)

var (
	ErrFileRequired  = errors.New("file is required (key: 'file')")
	ErrFileTooLarge  = errors.New("file size exceeds 5MB limit")
	ErrInvalidFormat = errors.New("invalid file format. Allowed formats: jpg, jpeg, png, webp, gif")
)

var AllowedMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format(layout)
	return []byte(`"` + formatted + `"`), nil
}

func (t *JsonTime) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(`"`+layout+`"`, string(data))
	if err != nil {
		return err
	}
	*t = JsonTime(parsed)
	return nil
}

func (t JsonTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *JsonTime) Scan(src any) error {
	if src == nil {
		*t = JsonTime(time.Time{})
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		*t = JsonTime(v)
		return nil
	case []byte:
		return t.parseString(string(v))
	case string:
		return t.parseString(v)
	default:
		return fmt.Errorf("cannot scan %T into JsonTime", src)
	}
}

func (t *JsonTime) parseString(str string) error {
	layouts := []string{
		layout,
		time.RFC3339,
		"2006-01-02 15:04:05.000000",
		"2006-01-02",
	}
	for _, l := range layouts {
		if parsed, err := time.Parse(l, str); err == nil {
			*t = JsonTime(parsed)
			return nil
		}
	}
	return fmt.Errorf("cannot parse %q into JsonTime", str)
}

func SetAccessTokenCookie(c *gin.Context, token string) {
	SetCookie(c, "access_token", token, 24*time.Hour)
}

func SetCookie(c *gin.Context, name, token string, maxExpired time.Duration) {
	isSecure := c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https"
	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		MaxAge:   int(maxExpired.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, cookie)
}

func ClearCookie(c *gin.Context, name string) {
	isSecure := c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https"
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, cookie)
}

func RemoveUploadFile(filenameOrPath string, uploadDirs ...string) error {
	trimmed := strings.TrimSpace(filenameOrPath)
	if trimmed == "" {
		return nil
	}

	if idx := strings.Index(trimmed, "?"); idx != -1 {
		trimmed = trimmed[:idx]
	}

	normalized := strings.ReplaceAll(trimmed, "\\", "/")
	filename := path.Base(normalized)

	if filename == "." || filename == "/" || filename == "" {
		return nil
	}

	uploadDir := DefaultUploadDir
	if len(uploadDirs) > 0 && uploadDirs[0] != "" {
		uploadDir = uploadDirs[0]
	}

	targetPath := filepath.Join(uploadDir, filename)
	if err := os.Remove(targetPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to remove upload file %s: %w", targetPath, err)
	}

	log.Printf("successfully removed physical upload file: %s", targetPath)
	return nil
}
