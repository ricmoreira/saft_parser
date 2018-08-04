package config

import (
	"os"
)

var (
	HOST              string = "HOST"
	MONGO_HOST        string = "MONGO_HOST"
	MONGO_DATABASE    string = "MONGO_DATABASE"
	SAFT_FILES_FOLDER string = "SAFT_FILES_FOLDER"
)

type Config struct {
	Host              string
	MongoHost         string
	MongoDatabaseName string
	SaftFilesFolder   string
}

func NewConfig() *Config {
	return &Config{
		Host:              MustGetEnv(HOST),
		MongoHost:         MustGetEnv(MONGO_HOST),
		MongoDatabaseName: MustGetEnv(MONGO_DATABASE),
		SaftFilesFolder:   MustGetEnv(SAFT_FILES_FOLDER),
	}
}

func MustGetEnv(envVarName string) string {
	res, found := os.LookupEnv(envVarName)

	if !found {
		panic("Environment variable " + envVarName + " not found")
	}

	return res
}
