package conf

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	AppUserName = "sourcescore"
	DbName      = "sourcescore"
)

type conf struct {
	AppUserPassword   string `env:"APP_USER_PASSWORD" yaml:"APP_USER_PASSWORD" env-required:"true"`
	PgHost            string `env:"PG_HOST" yaml:"PG_HOST" env-required:"true"`
	Port              string `env:"PORT" yaml:"PORT"`
	SuperUserPassword string `env:"SUPER_USER_PASSWORD" yaml:"SUPER_USER_PASSWORD" env-required:"true"`
}

var Cfg conf

func LoadConfig() {
	if envPath, ok := os.LookupEnv("DOTENV_PATH"); ok {
		file, err := os.Open(envPath)
		if err != nil {
			log.Fatalf("error while reading dotenv file: %s :: %s", envPath, err)
		}

		err = cleanenv.ParseYAML(file, &Cfg)
		if err != nil {
			log.Fatalf("error while parsing dotenv file: %s :: %s", envPath, err)
		}
	} else {
		err := cleanenv.ReadEnv(&Cfg)
		if err != nil {
			log.Fatalf("error while reading config environment variables :: %s", err)
		}
	}
}
