package config

type App struct {
	PORT          string `env:"APP_PORT" envDefault:"8080"`
	DB            Database
	Redis         RedisConfig
	MongoDbConfig MongoDbConfig
}
