package main

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Zzocker/book-labs/mediafile/config"
)

func main() {
	var cfg config.MediaFileServiceConfig
	{
		envFlag := flag.Bool("env", false, "true for reding all the configs from the environment variable")
		yamlFlag := flag.String("yml", "config/config.dev.yaml", "yaml config path")
		flag.Parse()

		var err error
		if *envFlag {
			err = cleanenv.ReadEnv(&cfg)
		} else {
			err = cleanenv.ReadConfig(*yamlFlag, &cfg)
		}
		if err != nil {
			panic(err)
		}
	}
	run(&cfg)
}
