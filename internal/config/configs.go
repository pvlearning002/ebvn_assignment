package config

import "github.com/kelseyhightower/envconfig"

/* General Config */

type Config struct {
	AppPort               int    `default:"8080" envconfig:"APP_PORT"`
	AppName               string `default:"MyApp" envconfig:"APP_NAME"`
	PasswordLengthDefault int    `default:"12" envconfig:"PASSWORD_LENGTH_DEFAULT"`
	Dash                  string `default:"-" envconfig:"DASH"`
	PasswordTitle         string `default:"Password" envconfig:"PASSWORD_TITLE"`
	StringEmpty           string `default:"" envconfig:"STRING_EMPTY"`
	StringSpace           string `default:" " envconfig:"STRING_SPACE"`
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

/* Status Config */

type ConfigStatus struct {
	StatusOK string `default:"OK" envconfig:"STATUS_OK"`
}

func NewConfigStatus() (*ConfigStatus, error) {
	cfg := &ConfigStatus{}
	err := envconfig.Process("config_status", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func CreateNewConfigStatus() *ConfigStatus {
	cfg, err := NewConfigStatus()
	if err != nil {
		panic(err)
	}
	return cfg
}

/* Status API Paths */

type ConfigPaths struct {
	PathGetGenPass     string `default:"/genpass" envconfig:"PATH_GET_GEN_PASS"`
	PathGetHealthCheck string `default:"/health" envconfig:"PATH_GET_HEALTH_CHECK"`
}

func NewConfigPaths() (*ConfigPaths, error) {
	cfg := &ConfigPaths{}
	err := envconfig.Process("config_paths", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func CreateNewConfigPaths() *ConfigPaths {
	cfg, err := NewConfigPaths()
	if err != nil {
		panic(err)
	}
	return cfg
}

/* Error Config */
type ConfigError struct {
	ErrorInternalServer string `default:"Internal Server Error" envconfig:"ERROR_INTERNAL_SERVER"`
	ErrorTitle          string `default:"Error" envconfig:"ERROR_TITLE"`
}

func NewConfigError() (*ConfigError, error) {
	cfg := &ConfigError{}
	err := envconfig.Process("config_error", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func CreateNewConfigError() *ConfigError {
	cfg, err := NewConfigError()
	if err != nil {
		panic(err)
	}
	return cfg
}

/* Health Check Config */

type ConfigHealthCheck struct {
	BookmarkService string `default:"bookmark_service" envconfig:"BOOKMARK_SERVICE"`
}

func NewConfigHealthCheck() (*ConfigHealthCheck, error) {
	cfg := &ConfigHealthCheck{}
	err := envconfig.Process("config_health_check", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func CreateNewConfigHealthCheck() *ConfigHealthCheck {
	cfg, err := NewConfigHealthCheck()
	if err != nil {
		panic(err)
	}
	return cfg
}
