package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rpc-ag/rpc-aggregator/internal/config"
	"github.com/rpc-ag/rpc-aggregator/internal/webserver"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()

	var configFilePath string
	var authFilePath string
	flag.StringVar(&configFilePath, "config", "./config.yaml", "config file to load")
	flag.StringVar(&authFilePath, "auth", "./auth.yaml", "auth file to load")
	flag.Parse()
	conf, err := config.ReadConfig(configFilePath)
	if err != nil {
		logger.Panic("failed to load config", zap.Error(err))
	}
	logger.Info("config loaded", zap.Any("conf", conf))

	auth, err := config.ReadAuth(authFilePath)
	if err != nil {
		logger.Panic("failed to load auth", zap.Error(err))
	}
	logger.Info("auth loaded")

	server, err := webserver.New(conf, auth, logger)
	if err != nil {
		logger.Panic("failed to start webserver", zap.Error(err))
	}

	logger.Info("Starting RPC Aggregator...")

	go func() {
		er := server.Run()
		if er != nil {
			panic(er)
		}
	}()

	go server.StartHealthChecker()

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)
	<-interrupt

	// Gracefully shutdown server
	logger.Info("Shutting down RPC Aggregator...")
	if err := server.Close(); err != nil {
		logger.Error("Error shutting down server")
	}
}
