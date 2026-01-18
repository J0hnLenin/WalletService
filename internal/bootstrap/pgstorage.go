package bootstrap

import (
	"fmt"

	"github.com/J0hnLenin/WalletService/config"
	"github.com/J0hnLenin/WalletService/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.StorageConfig) *pgstorage.PGStorage {

	storage, err := pgstorage.NewPGStorge(cfg)
	if err != nil {
		err = fmt.Errorf("init storage error: %w", err)
		panic(err)
	}
	return storage
}
