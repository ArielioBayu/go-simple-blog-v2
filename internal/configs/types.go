package configs

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
		DbSourceName string `mapstructure:"dbsourcename"`
	}

	SMTP struct {
		SMTPHost        string `mapstructure:"smtp_host"`
		SMTPPort        int    `mapstructure:"smtp_port"`
		SMTPSenderName  string `mapstructure:"smtp_sender_name"`
		SMTPSenderEmail string `mapstructure:"smtp_sender_email"`
		SMTPPassword    string `mapstructure:"smtp_password"`
	}
)
