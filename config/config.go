package config

import (
	"fmt"
	"os"
)

type MongoDb struct {
	Password string
	Name     string
}

func (m MongoDb) String() string {
	return fmt.Sprintf("name:%v\n", m.Name)
}

func MongoDbConfig() MongoDb {
	dbConfig := MongoDb{
		Password: os.Getenv("MONGO_DB_PASSWORD"),
		Name:     os.Getenv("MONGO_DB_NAME"),
	}

	return dbConfig
}

type AppEnvironment struct {
	AppEnv string
	Port   string
}

func AppConfig() AppEnvironment {
	appConfig := AppEnvironment{
		AppEnv: os.Getenv("APP_ENV"),
		Port:   os.Getenv("PORT"),
	}

	return appConfig
}
