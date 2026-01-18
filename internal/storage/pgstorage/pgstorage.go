package pgstorage

import (
	"context"
	"fmt"

	"github.com/J0hnLenin/WalletService/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGStorage struct {
	numberOfBuckets uint16
	shards []*pgShard
}

type pgShard struct {
	db *pgxpool.Pool
}

func NewPGStorge(cfg *config.Config) (*PGStorage, error) {

	storage := &PGStorage{
		numberOfBuckets: cfg.NumberOfBuckets,
		shards: make([]*pgShard, len(cfg.StorageShards)),
	}
	
	for shardIndex, shardConfig := range cfg.StorageShards {
		var err error

		connectionString := fmt.Sprintf("postgres_%d://%s:%s@%s:%d/%s",
			shardIndex,
			shardConfig.Username, 
			shardConfig.Password, 
			shardConfig.Host, 
			shardConfig.Port, 
			shardConfig.DBName)
	
		storage.shards[shardIndex], err = newPGShard(connectionString)

		if err != nil {
			return nil, fmt.Errorf("can't init shard %d: %w", shardIndex, err)
		}
	}

	return storage, nil
}

func newPGShard(connString string) (*pgShard, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("can't parse connection string '%s': %w", connString, err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("connection to '%s' failed: %w", connString, err)
	}

	shard := &pgShard{
		db: db,
	}

	return shard, nil
}

