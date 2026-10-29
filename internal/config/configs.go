package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AppPort               int    `default:"8080" envconfig:"APP_PORT"`
	AppName               string `default:"MyApp" envconfig:"APP_NAME"`
	StatusOK              string `default:"OK" envconfig:"STATUS_OK"`
	BookmarkService       string `default:"bookmark_service" envconfig:"BOOKMARK_SERVICE"`
	PathGetGenPass        string `default:"/genpass" envconfig:"PATH_GET_GEN_PASS"`
	PathGetHealthCheck    string `default:"/health" envconfig:"PATH_GET_HEALTH_CHECK"`
	PasswordLengthDefault int    `default:"12" envconfig:"PASSWORD_LENGTH_DEFAULT"`
	Dash                  string `default:"-" envconfig:"DASH"`
	ErrorInternalServer   string `default:"Internal Server Error" envconfig:"ERROR_INTERNAL_SERVER"`
	ErrorTitle            string `default:"Error" envconfig:"ERROR_TITLE"`
	PasswordTitle         string `default:"Password" envconfig:"PASSWORD_TITLE"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("config", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func CreateNewConfig() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
