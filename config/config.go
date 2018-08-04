package config

import (
	"os"
	"strconv"
)

var (
	HOST              string = "HOST"
	MONGO_HOST        string = "MONGO_HOST"
	MONGO_DATABASE    string = "MONGO_DATABASE"
	SAFT_FILES_FOLDER string = "SAFT_FILES_FOLDER"
	BOOTSTRAP_SERVERS string = "BOOTSTRAP_SERVERS"
	MESSAGE_MAX_BYTES string = "MESSAGE_MAX_BYTES"
)

type Config struct {
	Host              string
	MongoHost         string
	MongoDatabaseName string
	SaftFilesFolder   string
	BootstrapServers  string // kafka brokers endpoints (separated by ",")
	MessageMaxBytes   int    // kafka max message size to send from producer
}

func NewConfig() *Config {

	i, _ := strconv.Atoi(MustGetEnv(MESSAGE_MAX_BYTES))
	
	return &Config{
		Host:              MustGetEnv(HOST),
		MongoHost:         MustGetEnv(MONGO_HOST),
		MongoDatabaseName: MustGetEnv(MONGO_DATABASE),
		SaftFilesFolder:   MustGetEnv(SAFT_FILES_FOLDER),
		BootstrapServers:  MustGetEnv(BOOTSTRAP_SERVERS),
		MessageMaxBytes:   i,
	}
}

func MustGetEnv(envVarName string) string {
	res, found := os.LookupEnv(envVarName)

	if !found {
		panic("Environment variable " + envVarName + " not found")
	}

	return res
}
