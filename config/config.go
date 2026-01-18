package config

type Config struct {
	StorageConfig *StorageConfig
	APIConfig     *APIConfig
}

type StorageConfig struct {
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

type APIConfig struct {
	GrpcPort       uint16 `env:"API_GRPC_PORT"`
	ApiGatewayPort uint16 `env:"API_GATEWAY_PORT"`
}

func LoadConfig() (*Config, error) {

	storageConfig, err := newStorageConfig()
	if err != nil {
		return nil, err
	}

	apiConfig, err := apiConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
		APIConfig:     apiConfig,
		StorageConfig: storageConfig,
	}

	return config, nil
}
