package main

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Zzocker/book-labs/config"
	"github.com/Zzocker/book-labs/internal/app"
	"github.com/Zzocker/book-labs/pkg/logger"
)

func main() {
	var cfg config.BookSharing
	{
		envFag := flag.Bool("env", false, "true for reading configs from environment variables")
		cfgYmlPath := flag.String("yml", "config/config.dev.yml", "yaml config path")
		flag.Parse()

		var err error
		if *envFag {
			err = cleanenv.ReadEnv(&cfg)
		} else {
			err = cleanenv.ReadConfig(*cfgYmlPath, &cfg)
		}
		if err != nil {
			panic(err)
		}
	}
	// setup logger
	logger.NewServiceLogger(cfg.App.LogLevel, "application", cfg.App.Version)
	app.Run(&cfg)
}
