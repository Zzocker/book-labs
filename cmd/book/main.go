package main

import (
	"flag"

	"github.com/Zzocker/book-labs/config"
	"github.com/Zzocker/book-labs/internal/book"
	"github.com/Zzocker/book-labs/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
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
	logger.NewServiceLogger(cfg.App.LogLevel, "book", cfg.App.Version)
	book.Run(&cfg)
}
