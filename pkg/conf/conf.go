package conf

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"gorm.io/gorm"
)

const (
	AppUserName = "sourcescore"
)

type conf struct {
	AppUserPassword   string `env:"APP_USER_PASSWORD" yaml:"APP_USER_PASSWORD" env-required:"true"`
	DbName            string `env:"DB_NAME" yaml:"DB_NAME" env-default:"sourcescore"`
	PgHost            string `env:"PG_HOST" yaml:"PG_HOST" env-required:"true"`
	Port              string `env:"PORT" yaml:"PORT" env-default:"8080"`
	SuperUserPassword string `env:"SUPER_USER_PASSWORD" yaml:"SUPER_USER_PASSWORD" env-required:"true"`
	JwtSecret         string `env:"JWT_SECRET" yaml:"JWT_SECRET" env-required:"true"`
}

var (
	Cfg conf

	GormConfig = &gorm.Config{
		TranslateError: true,
	}
)

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
