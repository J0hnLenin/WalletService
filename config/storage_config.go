package config

import "fmt"

func newStorageConfig() (*StorageConfig, error) {
	buckets, err := getEnvInt("BUCKETS")
	if err != nil {
		return nil, err
	}

	shards, err := getEnvInt("SHARDS")
	if err != nil {
		return nil, err
	}

	storageConfig := &StorageConfig{
		NumberOfBuckets: uint16(buckets),
		StorageShards:   make([]*ShardConfig, shards),
	}

	for shardIndex := range shards {
		storageConfig.StorageShards[shardIndex], err = newShardConfig(shardIndex)
		if err != nil {
			return nil, fmt.Errorf("can't create storage shard %d: %w", shardIndex, err)
		}
	}

	return storageConfig, nil
}

func newShardConfig(shardIndex int) (*ShardConfig, error) {
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
		Host:     host,
		Port:     uint16(port),
		Username: username,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}, nil
}
