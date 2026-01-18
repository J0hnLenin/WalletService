package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/J0hnLenin/WalletService/config"
	server "github.com/J0hnLenin/WalletService/internal/api/wallet_service_api"
	"github.com/J0hnLenin/WalletService/internal/pb/wallets_api"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(cfg *config.APIConfig, api server.WalletServiceAPI) {
	go func() {
		if err := runGRPCServer(cfg, &api); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %w", err))
		}
	}()

	if err := runGatewayServer(cfg); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %w", err))
	}
}

func runGRPCServer(cfg *config.APIConfig, api *server.WalletServiceAPI) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	wallets_api.RegisterWalletsServiceServer(s, api)

	logString := fmt.Sprintf("gRPC-server server listening on :%d", cfg.GrpcPort)
	slog.Info(logString)
	return s.Serve(lis)
}

func runGatewayServer(cfg *config.APIConfig) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := "./internal/pb/swagger/wallets_api/wallets.swagger.json"
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		panic(fmt.Errorf("swagger file not found: %s", swaggerPath))
	}

	r := chi.NewRouter()
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	port := fmt.Sprintf(":%d", cfg.GrpcPort)
	err := wallets_api.RegisterWalletsServiceHandlerFromEndpoint(ctx, mux, port, opts)
	if err != nil {
		panic(err)
	}

	r.Mount("/", mux)
	logString := fmt.Sprintf("gRPC-Gateway server listening on :%d", cfg.ApiGatewayPort)
	slog.Info(logString)
	port = fmt.Sprintf(":%d", cfg.ApiGatewayPort)
	return http.ListenAndServe(port, r)
}
