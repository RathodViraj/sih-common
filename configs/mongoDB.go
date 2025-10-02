package configs

import "os"

type DBConfig struct {
	MongoURI string
	Database string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		MongoURI: os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DB"),
	}
}
