package configuration

import (
	env "github.com/Netflix/go-env"
)

type Config struct {
	HttpSocket string `env:"HTTP_SOCKET"`
	Postgresql Postgresql
}

type Postgresql struct {
	Host string `env:"DB_HOST"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
	Name string `env:"DB_NAME"`
}

func Load() (*Config, error) {
	var output Config
	_, err := env.UnmarshalFromEnviron(&output)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
