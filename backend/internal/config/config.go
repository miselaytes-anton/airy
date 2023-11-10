package config

import "os"

func GetBrokerAdress() string {
	value, ok := os.LookupEnv("BROKER_ADDRESS")
	if !ok {
		panic("BROKER_ADDRESS environment variable not set")
	}
	return value
}

func GetPostgresAddress() string {
	value, ok := os.LookupEnv("POSTGRES_ADDRESS")
	if !ok {
		panic("POSTGRES_ADDRESS environment variable not set")
	}
	return value
}
