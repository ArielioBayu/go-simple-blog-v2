package configs

import (
	"fmt"
	"net/url"
	"os"
)

type (
	Config struct {
		Service  Service      `mapstructure:"service"`
		Database Database     `mapstructure:"database"`
		SMTP     SMTP         `mapstructure:"smtp"`
		OAuth    OAuth        `mapstructure:"oauth"`
		Storage  Storage      `mapstructure:"storage"`
		R2       CloudflareR2 `mapstructure:"r2"`
	}
	Service struct {
		SecretKey   string `mapstructure:"secret_key"`
		Port        string `mapstructure:"port"`
		AppURL      string `mapstructure:"app_url"`
		FrontendURL string `mapstructure:"frontend_url"`
	}

	Database struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		Name            string `mapstructure:"name"`
		Timezone        string `mapstructure:"timezone"`
		ParseTime       bool   `mapstructure:"parse_time"`
		MultiStatements bool   `mapstructure:"multi_statements"`
	}

	SMTP struct {
		SMTPHost        string `mapstructure:"smtp_host"`
		SMTPPort        int    `mapstructure:"smtp_port"`
		SMTPSenderName  string `mapstructure:"smtp_sender_name"`
		SMTPSenderEmail string `mapstructure:"smtp_sender_email"`
		SMTPPassword    string `mapstructure:"smtp_password"`
	}

	OAuth struct {
		Google GoogleOAuth `mapstructure:"google"`
	}

	GoogleOAuth struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	}

	Storage struct {
		CloudflareR2 CloudflareR2 `mapstructure:"cloudflare_R2"`
	}

	CloudflareR2 struct {
		AccountID       string `mapstructure:"account_id"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		BucketName      string `mapstructure:"bucket_name"`
		PublicURL       string `mapstructure:"public_url"`

		// Alias untuk fleksibilitas konfigurasi yaml kapital
		R2AccountID       string `mapstructure:"R2_ACCOUNT_ID"`
		R2AccessKeyID     string `mapstructure:"R2_ACCESS_KEY_ID"`
		R2SecretAccessKey string `mapstructure:"R2_SECRET_ACCESS_KEY"`
		R2BucketName      string `mapstructure:"R2_BUCKET_NAME"`
		R2PublicURL       string `mapstructure:"R2_PUBLIC_URL"`
	}
)

// GetR2Config menyelesaikan konfigurasi Cloudflare R2 secara cerdas dan toleran
func (c *Config) GetR2Config() CloudflareR2 {
	res := c.R2

	// Prioritaskan storage.cloudflare_R2 jika ada
	if c.Storage.CloudflareR2.AccountID != "" {
		res.AccountID = c.Storage.CloudflareR2.AccountID
	}
	if c.Storage.CloudflareR2.R2AccountID != "" {
		res.AccountID = c.Storage.CloudflareR2.R2AccountID
	}
	if c.Storage.CloudflareR2.AccessKeyID != "" {
		res.AccessKeyID = c.Storage.CloudflareR2.AccessKeyID
	}
	if c.Storage.CloudflareR2.R2AccessKeyID != "" {
		res.AccessKeyID = c.Storage.CloudflareR2.R2AccessKeyID
	}
	if c.Storage.CloudflareR2.SecretAccessKey != "" {
		res.SecretAccessKey = c.Storage.CloudflareR2.SecretAccessKey
	}
	if c.Storage.CloudflareR2.R2SecretAccessKey != "" {
		res.SecretAccessKey = c.Storage.CloudflareR2.R2SecretAccessKey
	}
	if c.Storage.CloudflareR2.BucketName != "" {
		res.BucketName = c.Storage.CloudflareR2.BucketName
	}
	if c.Storage.CloudflareR2.R2BucketName != "" {
		res.BucketName = c.Storage.CloudflareR2.R2BucketName
	}
	if c.Storage.CloudflareR2.PublicURL != "" {
		res.PublicURL = c.Storage.CloudflareR2.PublicURL
	}
	if c.Storage.CloudflareR2.R2PublicURL != "" {
		res.PublicURL = c.Storage.CloudflareR2.R2PublicURL
	}

	// Fallback ke environment variables jika belum diisi di file
	if res.AccountID == "" {
		res.AccountID = os.Getenv("R2_ACCOUNT_ID")
	}
	if res.AccessKeyID == "" {
		res.AccessKeyID = os.Getenv("R2_ACCESS_KEY_ID")
	}
	if res.SecretAccessKey == "" {
		res.SecretAccessKey = os.Getenv("R2_SECRET_ACCESS_KEY")
	}
	if res.BucketName == "" {
		res.BucketName = os.Getenv("R2_BUCKET_NAME")
	}
	if res.PublicURL == "" {
		res.PublicURL = os.Getenv("R2_PUBLIC_URL")
	}

	return res
}

// GetDSN merakit string koneksi database secara dinamis dan aman
func (d Database) GetDSN() string {
	tz := d.Timezone
	if tz == "" {
		tz = "Asia/Jakarta"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=%t&loc=%s&multiStatements=%t",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.ParseTime,
		url.QueryEscape(tz),
		d.MultiStatements,
	)
}
