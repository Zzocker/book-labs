package config

type MediaFileServiceConfig struct {
	Log      Log      `yaml:"log"`
	S3Config S3Config `yaml:"s3_config"`
	Port     string   `yaml:"port" env:"MEDIAFILE_PORT" env-required:"true"`
}

type Log struct {
	Level          string `yaml:"level" env:"MEDIAFILE_LOG_LEVEL" env-required:"true" env-default:"info"`
	ServiceName    string `yaml:"name" env:"MEDIAFILE_LOG_NAME" env-required:"true" env-default:"mediafile"`
	ServiceVersion string `yaml:"version" env:"MEDIAFILE_LOG_VERSION" env-required:"true" env-default:"1.0.0"`
}

type S3Config struct {
	Endpoint        string `yaml:"endpoint" env:"MEDIAFILE_S3CONFIG_ENDPOINT" env-required:"true"`
	AccessKeyID     string `yaml:"accessKeyID" env:"MEDIAFILE_S3CONFIG_ACCESSKEY_ID" env-required:"true"`
	SecretAccessKey string `yaml:"secretAccessKey" env:"MEDIAFILE_S3CONFIG_SECRET_ACCESS_KEY" env-required:"true"`
	SessionToken    string `yaml:"sessionToken" env:"MEDIAFILE_S3CONFIG_SESSION_TOKEN" env-default:""`
	Region          string `yaml:"region" env:"MEDIAFILE_S3CONFIG_REGION" env-required:"true"`
	BucketName      string `yaml:"bucketName" env:"MEDIAFILE_S3CONFIG_BUCKET_NAME" env-required:"true"`
}
