package config

type BookSharing struct {
	App App `yaml:"app"`
}

type App struct {
	Port       string     `yaml:"port" env:"APP_PORT" env-required:"true"`
	Version    string     `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	LogLevel   string     `yaml:"log_level" env:"APP_LOG_LEVEL" env-required:"true"`
	OAuthRedis OAuthRedis `yaml:"oauth_redis"`
}

type OAuthRedis struct {
	ExpiryS  int64  `yaml:"expiry_s" env:"APP_OAUTH_REDIS_EXPIRY_S" env-default:"604800"` // 7 days
	URL      string `yaml:"url" env:"APP_OAUTH_REDIS_URL" env-required:"true"`
	Username string `yaml:"username" env:"APP_OAUTH_REDIS_USERNAME" env-default:""`
	Password string `yaml:"password" env:"APP_OAUTH_REDIS_PASSWORD" env-default:""`
	Database int    `yaml:"database" env:"APP_OAUTH_REDIS_DATABASE" env-required:"true"`
}
