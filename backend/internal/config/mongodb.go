package config

type MongoDbConfig struct {
	Addr         string `env:"MONGODB_ADDRESS"`
	DatabaseName string `env:"MONGODB_DATABASE_NAME"`
}
