package config

import "github.com/kelseyhightower/envconfig"

/* General Config */

type Config struct {
	AppPort               int    `default:"8080" envconfig:"APP_PORT"`
	AppName               string `default:"MyApp" envconfig:"APP_NAME"`
	PasswordLengthDefault int    `default:"12" envconfig:"PASSWORD_LENGTH_DEFAULT"`
	KeycodeLengthDefault  int    `default:"7" envconfig:"KEYCODE_LENGTH_DEFAULT"`
	Dash                  string `default:"-" envconfig:"DASH"`
	PasswordTitle         string `default:"Password" envconfig:"PASSWORD_TITLE"`
	KeyTitle              string `default:"key" envconfig:"KEY_TITLE"`
	StringEmpty           string `default:"" envconfig:"STRING_EMPTY"`
	StringSpace           string `default:" " envconfig:"STRING_SPACE"`
	RedisExpiredTime      int    `default:"3600" envconfig:"REDIS_EXPIRED_TIME"`
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
	PathSwagger        string `default:"/swagger/*any" envconfig:"PATH_GET_GEN_PASS_SWAGGER"`
	PathGetGenPass     string `default:"/genpass" envconfig:"PATH_GET_GEN_PASS"`
	PathGetHealthCheck string `default:"/health" envconfig:"PATH_GET_HEALTH_CHECK"`
	PathPostShortenURL string `default:"/v1/links/shorten" envconfig:"PATH_POST_SHORTEN_URL"`
	PathGetURL         string `default:"/v1/links/redirect/:code" envconfig:"PATH_GET_URL"`
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
	ErrorTitle           string `default:"Error" envconfig:"ERROR_TITLE"`
	ErrorInternalServer  string `default:"Internal Server Error" envconfig:"ERROR_INTERNAL_SERVER"`
	ErrorInvalidInput    string `default:"Invalid input" envconfig:"ERROR_INVALID_INPUT"`
	ErrorRedisConnection string `default:"Redis connection error" envconfig:"ERROR_REDIS_CONNECTION"`
	ErrorNotFound        string `default:"Not Found" envconfig:"ERROR_NOT_FOUND"`
	ErrorLengthZero      string `default:"Length of string must be greater than 0" envconfig:"ERROR_LENGTH_ZERO"`
	ErrorInvalidCode     string `default:"Invalid code" envconfig:"ERROR_INVALID_CODE"`
	ErrorUrlNotFound     string `default:"URL not found" envconfig:"ERROR_URL_NOT_FOUND"`
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

/* Parameter Config */
type ConfigParameter struct {
	Code string `default:"code" envconfig:"CODE"`
}

func NewConfigParameter() (*ConfigParameter, error) {
	cfg := &ConfigParameter{}
	err := envconfig.Process("config_parameter", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func CreateNewConfigParameter() *ConfigParameter {
	cfg, err := NewConfigParameter()
	if err != nil {
		panic(err)
	}
	return cfg
}
