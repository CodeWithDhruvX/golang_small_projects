package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	MongoHost        string
	MongoPort        string
	MongoUser        string
	MongoPassword    string
	MongoDB          string
}

type AppConfig struct {
	ServerPort string
	Database   DatabaseConfig
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		ServerPort: getEnv("SERVER_PORT", "8082"),
		Database: DatabaseConfig{
			PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
			PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
			PostgresUser:     getEnv("POSTGRES_USER", "admin"),
			PostgresPassword: getEnv("POSTGRES_PASSWORD", "password123"),
			PostgresDB:       getEnv("POSTGRES_DB", "project7_db"),
			MongoHost:        getEnv("MONGO_HOST", "localhost"),
			MongoPort:        getEnv("MONGO_PORT", "27017"),
			MongoUser:        getEnv("MONGO_USER", "admin"),
			MongoPassword:    getEnv("MONGO_PASSWORD", "password123"),
			MongoDB:          getEnv("MONGO_DB", "project7_mongo"),
		},
	}
}

func (c *AppConfig) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.Database.PostgresHost,
		c.Database.PostgresUser,
		c.Database.PostgresPassword,
		c.Database.PostgresDB,
		c.Database.PostgresPort,
	)
}

func (c *AppConfig) GetMongoURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s",
		c.Database.MongoUser,
		c.Database.MongoPassword,
		c.Database.MongoHost,
		c.Database.MongoPort,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
