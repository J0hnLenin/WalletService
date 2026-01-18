package config

import (
	"fmt"
	"os"
	"strconv"
)
type Config struct {
	NumberOfBuckets uint16 `env:"BUCKETS"`
	StorageShards   []*StorageConfig
}

type StorageConfig struct {
	Host     string `env:"STORAGE_HOST_N"`
	Port     uint16 `env:"STORAGE_PORT_N"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB"`
	SSLMode  string `env:"STORAGE_SSL_MODE"`
}

func LoadConfig() (*Config, error) {
	buckets, err := strconv.Atoi(os.Getenv("BUCKETS"))
	if err != nil {
		return nil, err
	}

	shards, err := strconv.Atoi(os.Getenv("SHARDS"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		NumberOfBuckets: uint16(buckets),
		StorageShards: make([]*StorageConfig, shards),
	}

	for shardIndex := range(shards) {
		config.StorageShards[shardIndex], err = newStorageConfig(shardIndex)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func newStorageConfig(shardIndex int) (*StorageConfig, error) {
	
	hostEnv := fmt.Sprintf("STORAGE_HOST_%d", shardIndex)
	host := os.Getenv(hostEnv)
	
	portEnv := fmt.Sprintf("STORAGE_PORT_%d", shardIndex)
	port, err := strconv.Atoi(os.Getenv(portEnv))
	if err != nil {
		return nil, err
	}

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("STORAGE_SSL_MODE")

	if host == "" ||
		username == "" ||
		password == "" ||
		dbName == "" ||
		sslMode == "" {

		return nil, fmt.Errorf("Storage config for shard %d not found", shardIndex)
	}

	return &StorageConfig{
		Host: host,
		Port: uint16(port),
		Username: username,
		Password: password,
		DBName: dbName,
		SSLMode: sslMode,
	}, nil
	
}