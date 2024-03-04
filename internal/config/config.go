package config

type Config struct {
	ServicePort string `env:"SERVICE_PORT,required=true"`
}
