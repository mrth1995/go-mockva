package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/caarlos0/env"
)

type Config struct {
	Port             int    `env:"PORT" envDocs:"Application port" envDefault:"8080"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDocs:"PostgreSQL port"`
	PostgresHost     string `env:"POSTGRES_HOST" envDocs:"PostgreSQL host"`
	PostgresUsername string `env:"POSTGRES_USERNAME" envDocs:"PostgreSQL username"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDocs:"PostgreSQL password"`
	DBName           string `env:"DB_NAME" envDocs:"Database name" envDefault:"mockva"`
	SQLFilePath      string `env:"SQL_FILE_PATH" envDocs:"SQL file path for schema migration" envDefault:"/srv/migration"`
	SwaggerFilePath  string `env:"SWAGGER_FILE_PATH"`
}

func (envVar Config) HelpDocs() []string {
	reflectEnvVar := reflect.TypeOf(envVar)
	doc := make([]string, 1+reflectEnvVar.NumField())
	doc[0] = "Environment variables config:"
	for i := 0; i < reflectEnvVar.NumField(); i++ {
		field := reflectEnvVar.Field(i)
		envName := field.Tag.Get("env")
		envDefault := field.Tag.Get("envDefault")
		envDocs := field.Tag.Get("envDocs")
		doc[i+1] = fmt.Sprintf("  %v\t %v (default: %v)", envName, envDocs, envDefault)
	}
	return doc
}

func ParseConfiguration() (*Config, error) {
	cfg := &Config{}
	flag.Usage = func() {
		flag.CommandLine.SetOutput(os.Stdout)
		for _, val := range cfg.HelpDocs() {
			fmt.Println(val)
		}
		fmt.Println("")
		flag.PrintDefaults()
	}
	flag.Parse()
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
