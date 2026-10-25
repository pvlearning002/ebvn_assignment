package config

var cfg *Config = CreateNewConfig()

func InitConfig() *Config {
	if cfg == nil {
		cfg = CreateNewConfig()
	}
	return cfg
}

var (
	APP_PORT                = cfg.AppPort
	OK_STATUS               = cfg.StatusOK
	BOOKMARK_SERVICE        = cfg.BookmarkService
	PATH_GET_HEALTH_CHECK   = cfg.PathGetHealthCheck
	PATH_GET_GEN_PASS       = cfg.PathGetGenPass
	PASSWORD_LENGTH_DEFAULT = cfg.PasswordLengthDefault
)
