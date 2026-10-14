package constants

import (
	"fmt"
	"strings"
	"time"
)

const (
	LayoutDateTime           = "2006-01-02 15:04:05"
	OTPTypeEmailVerification = "EMAIL_VERIFICATION"
	OTPExpirationDuration    = 5 * time.Minute
	DefaultUploadDir         = "./uploads"
	UploadDir                = "./uploads"
	MaxUploadSize            = 5 * 1024 * 1024 // 5MB

	// Follow Statuses
	FollowStatusPending  = "pending"
	FollowStatusAccepted = "accepted"

	// Notification Messages
	MsgLikePost      = "menyukai postingan Anda"
	MsgFollow        = "mulai mengikuti Anda"
	MsgFollowRequest = "meminta untuk mengikuti Anda"
)

func FormatCommentNotification(content string) string {
	preview := strings.TrimSpace(content)
	runes := []rune(preview)
	if len(runes) > 50 {
		preview = string(runes[:47]) + "..."
	}
	return fmt.Sprintf("mengomentari postingan Anda: %s", preview)
}

func FormatReplyNotification(content string) string {
	preview := strings.TrimSpace(content)
	runes := []rune(preview)
	if len(runes) > 50 {
		preview = string(runes[:47]) + "..."
	}
	return fmt.Sprintf("membalas komentar Anda: \"%s\"", preview)
}
