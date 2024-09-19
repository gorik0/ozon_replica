package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

func MustLoad() *Config {

	var cfg Config
	//	::: read env
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Printf("Error loading config: %s", err)
		os.Exit(1)
	}
	//var grps GRPC
	//err = cleanenv.ReadEnv(&grps)
	//log.Printf("GRPC :::  %v\n", grps)
	//if err != nil {
	//	log.Printf("Error loading config grpc: %s", err)
	//	os.Exit(1)
	//}
	//cfg.GRPC = grps
	//	::: read cfg
	if _, err = os.Stat(cfg.ConfigPath); os.IsNotExist(err) {

		log.Printf("Error reading path to config: %s", err)
		os.Exit(1)

	}
	if err != nil {
		return nil
	}

	err = cleanenv.ReadConfig(cfg.ConfigPath, &cfg)
	if err != nil {
		log.Printf("Error loading config from env: %s : %s", cfg.ConfigPath, err)
		os.Exit(1)
	}
	return &cfg
}
