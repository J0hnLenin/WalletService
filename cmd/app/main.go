package main

import (
	"log/slog"
	"os"

	"github.com/J0hnLenin/WalletService/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error("Can't load application config", 
			"err", err)
		panic(err)
	}
	logger.Info("Config loaded", "shards", len(config.StorageShards))
	
	logger.Info("Application started")
}