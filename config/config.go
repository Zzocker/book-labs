package config

type BookSharing struct {
	App         App         `yaml:"app"`
	UserProfile UserProfile `yaml:"userprofile"`
	MongoDB     MongoDB     `yaml:"mongo_db"`
	Book        Book        `yaml:"book"`
	S3          S3          `yaml:"s3"`
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

type UserProfile struct {
	Port              string `yaml:"port" env:"USERPROFILE_PORT" env-required:"true"`
	Version           string `yaml:"version" env:"USERPROFILE_VERSION" env-default:"1.0.0"`
	LogLevel          string `yaml:"log_level" env:"USERPROFILE_LOG_LEVEL" env-required:"true"`
	CollectionName    string `yaml:"collection_name" env:"USERPROFILE_COLLECTION_NAME" env-required:"true"`
	ProfileBucketName string `yaml:"s3_bucket_name" env:"USERPROFILE_S3_BUCKET" env-required:"true"`
}

type Book struct {
	Port           string `yaml:"port" env:"BOOK_PORT" env-required:"true"`
	Version        string `yaml:"version" env:"BOOK_VERSION" env-default:"1.0.0"`
	LogLevel       string `yaml:"log_level" env:"BOOK_LOG_LEVEL" env-required:"true"`
	CollectionName string `yaml:"collection_name" env:"BOOK_COLLECTION_NAME" env-required:"true"`
	BookBucketName string `yaml:"s3_bucket_name" env:"BOOK_S3_BUCKET" env-required:"true"`
}

type MongoDB struct {
	URL      string `yaml:"url" env:"MONGODB_URL" env-required:"true"`
	Username string `yaml:"username" env:"MONGODB_Username" env-required:"true"`
	Password string `yaml:"password" env:"MONGODB_Password" env-required:"true"`
	Database string `yaml:"database" env:"MONGODB_Database" env-required:"true"`
}

type S3 struct {
	Endpoint        string `yaml:"endpoint" env:"S3_ENDPOINT" env-required:"true"`
	AccessKeyID     string `yaml:"access_key_id" env:"S3_ACCESS_KEY_ID" env-required:"true"`
	SecretAccessKey string `yaml:"secret_access_key" env:"S3_SECRET_ACCESS_KEY" env-required:"true"`
	SessionToken    string `yaml:"session_token" env:"S3_SESSION_TOKE" env-default:""`
	Region          string `yaml:"region" env:"S3_REGION" env-required:"true"`
}
