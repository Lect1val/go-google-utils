package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Gcp gcpClient
}

type gcpClient struct {
	Secret string `env:"GCP_CLIENT_SECRET"`
	Token  string `env:"GCP_CLIENT_TOKEN"`
}

var Val Config
var once sync.Once

func prefix(e string) string {
	if e == "" {
		return ""
	}
	return fmt.Sprintf("%s_", e)
}

func C() Config {
	once.Do(func() {
		opts := env.Options{
			Prefix: prefix(Env),
		}

		conf, err := parseEnv[Config](opts)
		if err != nil {
			log.Fatal(err)
		}

		Val = conf
	})

	return Val
}
