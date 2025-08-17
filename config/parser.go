package config

import "github.com/caarlos0/env/v11"

func parseEnv[T any](opts env.Options) (T, error) {
	var t T
	if err := env.Parse(&t); err != nil {
		return t, err
	}

	if err := env.ParseWithOptions(&t, opts); err != nil {
		return t, err
	}
	return t, nil
}
