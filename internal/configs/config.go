package configs

import "github.com/spf13/viper"

var config *Config

type option struct {
	configFolders []string
	configFile    string
	configType    string
}

type Opsi func(*option)

func Init(opts ...Opsi) error {
	opt := &option{
		configFolders: getDefaultConfigFolder(),
		configFile:    getDefaultConfigFile(),
		configType:    getDefaultConfigType(),
	}

	for _, optfunc := range opts {
		optfunc(opt)
	}

	for _, ConfigFolder := range opt.configFolders {
		viper.AddConfigPath(ConfigFolder)
	}

	viper.SetConfigName(opt.configFile)
	viper.SetConfigType(opt.configType)
	viper.AutomaticEnv()

	//	untuk unmarshal file yaml kedalam variabel cfg
	config = new(Config)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&config)
}

func getDefaultConfigFolder() []string {
	return []string{"./configs/"}
}

func getDefaultConfigFile() string {
	return "config"
}

func getDefaultConfigType() string {
	return "yaml"
}

func WithConfigFolder(configFolder []string) Opsi {
	return func(opt *option) {
		opt.configFolders = configFolder
	}
}

func WithConfigFile(configFile string) Opsi {
	return func(opt *option) {
		opt.configFile = configFile
	}
}

func WithConfigType(configType string) Opsi {
	return func(opt *option) {
		opt.configType = configType
	}
}

func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}
