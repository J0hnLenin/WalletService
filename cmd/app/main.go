package main

import (
	"log/slog"
	"os"

	"github.com/J0hnLenin/WalletService/config"
	"github.com/J0hnLenin/WalletService/internal/bootstrap"
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
	
	walletStorage := bootstrap.InitPGStorage(config)
	walletService := bootstrap.InitWalletService(walletStorage)
	walletsApi := bootstrap.InitWalletServiceAPI(walletService)

	bootstrap.AppRun(*walletsApi)

	logger.Info("Application started")
}