package configs

type (
	Config struct {
		Service  Service  `mapstructure:"service"`
		Database Database `mapstructure:"database"`
	}
	Service struct {
		SecretKey string `mapstructure:"secret_key"`
		Port      string `mapstructure:"port"`
	}

	Database struct {
		DbSourceName string `mapstructure:"dbsourcename"`
	}
)
