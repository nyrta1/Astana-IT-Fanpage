package config

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDRESS"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}
