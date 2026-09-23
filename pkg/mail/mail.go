package mail

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
)

func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random otp: %w", err)
	}
	otp := n.Int64() + 100000
	return fmt.Sprintf("%06d", otp), nil
}

func SendOTPEmail(cfg configs.SMTP, toEmail, username, otpCode string) error {
	if cfg.SMTPSenderEmail == "" || cfg.SMTPPassword == "" {
		log.Printf("\n==================================================\n"+
			"📧 [EMAIL OTP DEV SIMULATION]\n"+
			"To: %s (%s)\n"+
			"OTP Code: %s\n"+
			"Expires in: 5 minutes\n"+
			"==================================================\n", toEmail, username, otpCode)
		return nil
	}

	senderName := cfg.SMTPSenderName
	if senderName == "" {
		senderName = "Simple Blog"
	}

	subject := fmt.Sprintf("%s - Kode Verifikasi Email OTP", otpCode)
	fromHeader := fmt.Sprintf("%s <%s>", senderName, cfg.SMTPSenderEmail)

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Verifikasi Email</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f1f5f9; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; color: #1e293b;">
  <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background: linear-gradient(180deg, #f8fafc 0%%, #e2e8f0 100%%); padding: 48px 16px;">
    <tr>
      <td align="center">
        <!-- Main Card -->
        <table role="presentation" width="100%%" style="max-width: 520px; background-color: #ffffff; border-radius: 20px; border: 1px solid #e2e8f0; overflow: hidden; box-shadow: 0 10px 30px -5px rgba(0, 0, 0, 0.06);">
          <!-- Top Gradient Accent Bar -->
          <tr>
            <td style="background: linear-gradient(90deg, #6366f1 0%%, #a855f7 50%%, #ec4899 100%%); height: 6px;"></td>
          </tr>
          
          <!-- Content Body -->
          <tr>
            <td style="padding: 40px 36px 36px 36px;">
              <!-- Brand Badge -->
              <table role="presentation" cellspacing="0" cellpadding="0" style="margin-bottom: 24px;">
                <tr>
                  <td style="background: #eef2ff; border: 1px solid #e0e7ff; border-radius: 20px; padding: 6px 14px; font-size: 13px; font-weight: 600; color: #4f46e5;">
                    ✨ %s
                  </td>
                </tr>
              </table>

              <!-- Greeting -->
              <h2 style="margin: 0 0 12px; font-size: 24px; font-weight: 700; color: #0f172a; letter-spacing: -0.5px;">
                Halo, @%s 👋
              </h2>
              <p style="margin: 0 0 28px; font-size: 15px; color: #475569; line-height: 1.6;">
                Terima kasih telah mendaftar di <strong>%s</strong>. Masukkan kode verifikasi berikut untuk menyelesaikan pendaftaran akun Anda:
              </p>
              
              <!-- OTP Box -->
              <div style="background: linear-gradient(135deg, #f8fafc 0%%, #eef2ff 100%%); border: 2px dashed #cbd5e1; border-radius: 16px; padding: 24px 16px; text-align: center; margin-bottom: 26px;">
                <span style="font-size: 36px; font-weight: 800; letter-spacing: 10px; color: #4338ca; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace; display: inline-block;">
                  %s
                </span>
              </div>
              
              <!-- Expiry Alert Pill -->
              <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="margin-bottom: 20px;">
                <tr>
                  <td style="background-color: #fffbeb; border: 1px solid #fef3c7; border-radius: 12px; padding: 12px 16px; font-size: 13px; color: #b45309; line-height: 1.5;">
                    ⏱️ Kode ini hanya berlaku selama <strong>5 menit</strong>. Jangan bagikan kode ini kepada siapa pun.
                  </td>
                </tr>
              </table>

              <!-- Security Disclaimer -->
              <p style="margin: 0 0 28px; font-size: 13px; color: #64748b; line-height: 1.5;">
                Jika Anda tidak merasa melakukan pendaftaran ini, abaikan email ini dengan aman. Akun tidak akan aktif tanpa verifikasi kode ini.
              </p>
              
              <!-- Divider -->
              <hr style="border: none; border-top: 1px solid #f1f5f9; margin: 0 0 24px 0;">
              
              <!-- Footer -->
              <p style="margin: 0; font-size: 12px; color: #94a3b8; text-align: center; line-height: 1.5;">
                &copy; %s. All rights reserved.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, senderName, username, senderName, otpCode, senderName)

	headers := make(map[string]string)
	headers["From"] = fromHeader
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n" + body)

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.SMTPSenderEmail, cfg.SMTPPassword, cfg.SMTPHost)

	err := smtp.SendMail(addr, auth, cfg.SMTPSenderEmail, []string{toEmail}, []byte(message.String()))
	if err != nil {
		return fmt.Errorf("failed to send otp email via smtp: %w", err)
	}

	return nil
}
