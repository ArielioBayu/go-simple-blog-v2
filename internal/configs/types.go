package configs

import (
	"fmt"
	"net/url"
)

type (
	Config struct {
		Service  Service  `mapstructure:"service"`
		Database Database `mapstructure:"database"`
		SMTP     SMTP     `mapstructure:"smtp"`
	}
	Service struct {
		SecretKey string `mapstructure:"secret_key"`
		Port      string `mapstructure:"port"`
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
)

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
