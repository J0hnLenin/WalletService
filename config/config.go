package config

import (
	"fmt"
	"os"
	"strconv"
)
type Config struct {
	NumberOfBuckets uint16 `env:"BUCKETS"`
	StorageShards   []*ShardConfig
}

type ShardConfig struct {
	Host     string `env:"STORAGE_HOST_N"`
	Port     uint16 `env:"STORAGE_PORT_N"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB"`
	SSLMode  string `env:"STORAGE_SSL_MODE"`
}

func LoadConfig() (*Config, error) {

	buckets, err := getEnvInt("BUCKETS")
	if err != nil {
		return nil, err
	}

	shards, err := getEnvInt("SHARDS")
	if err != nil {
		return nil, err
	}

	config := &Config{
		NumberOfBuckets: uint16(buckets),
		StorageShards: make([]*ShardConfig, shards),
	}

	for shardIndex := range(shards) {
		config.StorageShards[shardIndex], err = newStorageConfig(shardIndex)
		if err != nil {
			return nil, fmt.Errorf("can't create storage shard %d: %w", shardIndex, err)
		}
	}

	return config, nil
}

func newStorageConfig(shardIndex int) (*ShardConfig, error) {
	var envVariableName string

	envVariableName = fmt.Sprintf("STORAGE_HOST_%d", shardIndex)
	host, err := getEnv(envVariableName)
	if err != nil {
		return nil, err
	}

	envVariableName = fmt.Sprintf("STORAGE_PORT_%d", shardIndex)
	port, err := getEnvInt(envVariableName)
	if err != nil {
		return nil, err
	}

	username, err := getEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}

	password, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}

	sslMode, err := getEnv("STORAGE_SSL_MODE")
	if err != nil {
		return nil, err
	}

	return &ShardConfig{
		Host: host,
		Port: uint16(port),
		Username: username,
		Password: password,
		DBName: dbName,
		SSLMode: sslMode,
	}, nil
}

func getEnv(variableName string) (string, error) {
	value := os.Getenv(variableName)
	if value == "" {
		return value, fmt.Errorf("got empty value for '%s' from env", variableName)
	}
	return value, nil
}

func getEnvInt(variableName string) (int, error) {
	value := os.Getenv(variableName)
	result, err := strconv.Atoi(value)
	if err != nil {
		return result, fmt.Errorf("can't convert '%s' env value '%s' to int: %w", variableName, value, err)
	}
	return result, nil
}